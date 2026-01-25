package bridge

/*
int Process(const char* data, int len);
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"unicode/utf8"
	"unsafe"
)

//
// 1) Request types (meaningfully separated)
//

// シナリオ実行用：params 必須
type ScenarioExecuteReq struct {
	Type     string          `json:"type"`
	ID       int64           `json:"id"`
	HalType  string          `json:"hal_type"`
	Scenario string          `json:"scenario"`
	Params   json.RawMessage `json:"params"` // ★必須なのでomitempty無し
}

// シナリオ設定用：params 無し（必要なら Settings 等を追加してOK）
type ScenarioConfigReq struct {
	Type     string `json:"type"`
	ID       int64  `json:"id"`
	HalType  string `json:"hal_type"`
	Scenario string `json:"scenario"`
	// Paramsなし
}

//
// 2) Normalized payload (what we send to C++)
//

type Normalized struct {
	SchemaVersion int       `json:"schemaVersion"`
	RequestID     uint32    `json:"requestId"`
	Type          string    `json:"type"`
	ID            int64     `json:"id"`
	HalType       string    `json:"hal_type"`
	Scenario      string    `json:"scenario"`
	Params        []ParamNV `json:"params,omitempty"` // 設定用は省略される
}

// C++側を楽にするため value は文字列に統一（object/arrayはJSON文字列化）
type ParamNV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ParamKV struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

//
// 3) RequestId generator (wrap to 0)
//

type RequestIDGen struct {
	max uint32
	cur atomic.Uint32
}

func NewRequestIDGen(max uint32, start uint32) *RequestIDGen {
	g := &RequestIDGen{max: max}
	g.cur.Store(start)
	return g
}

func (g *RequestIDGen) Next() uint32 {
	for {
		old := g.cur.Load()
		var next uint32
		if old >= g.max {
			next = 0
		} else {
			next = old + 1
		}
		if g.cur.CompareAndSwap(old, next) {
			return next
		}
	}
}

//
// 4) Limits / policy
//

type Limits struct {
	MaxParams       int
	MaxKeyLen       int
	MaxValueLen     int
	MaxPayloadBytes int

	// object/array対策
	MaxDepth        int
	MaxContainerLen int
}

type DupKeyPolicy int

const (
	DupKeyError DupKeyPolicy = iota
	DupKeyLastWins
)

//
// 5) Public APIs
//

// 実行用：params必須、正規化→C++へ送る
func NormalizeSendExecute(raw []byte, lim Limits, gen *RequestIDGen, dup DupKeyPolicy) error {
	var req ScenarioExecuteReq
	if err := json.Unmarshal(raw, &req); err != nil {
		return fmt.Errorf("invalid execute json: %w", err)
	}

	payload, err := normalizeCore(coreIn{
		Type:          req.Type,
		ID:            req.ID,
		HalType:       req.HalType,
		Scenario:      req.Scenario,
		ParamsRaw:     req.Params,
		RequireParams: true,
	}, lim, gen, dup)
	if err != nil {
		return err
	}

	return sendToCpp(payload)
}

// 設定用：params無し、正規化→C++へ送る（※送らない運用なら sendToCpp を呼ばない設計にしてOK）
func NormalizeSendConfig(raw []byte, lim Limits, gen *RequestIDGen) error {
	var req ScenarioConfigReq
	if err := json.Unmarshal(raw, &req); err != nil {
		return fmt.Errorf("invalid config json: %w", err)
	}

	payload, err := normalizeCore(coreIn{
		Type:          req.Type,
		ID:            req.ID,
		HalType:       req.HalType,
		Scenario:      req.Scenario,
		ParamsRaw:     nil,
		RequireParams: false,
	}, lim, gen, DupKeyError)
	if err != nil {
		return err
	}

	return sendToCpp(payload)
}

//
// 6) Core normalize (shared)
//

type coreIn struct {
	Type          string
	ID            int64
	HalType       string
	Scenario      string
	ParamsRaw     json.RawMessage
	RequireParams bool
}

func normalizeCore(in coreIn, lim Limits, gen *RequestIDGen, dup DupKeyPolicy) ([]byte, error) {
	// top-level normalize
	in.Type = strings.TrimSpace(in.Type)
	in.HalType = strings.TrimSpace(in.HalType)
	in.Scenario = strings.TrimSpace(in.Scenario)

	// required checks
	if in.Type == "" || in.HalType == "" || in.Scenario == "" {
		return nil, errors.New("missing required top-level fields (type/hal_type/scenario)")
	}
	if in.ID <= 0 {
		return nil, fmt.Errorf("invalid id: %d", in.ID)
	}

	// requestId
	reqID := uint32(0)
	if gen != nil {
		reqID = gen.Next()
	}

	n := Normalized{
		SchemaVersion: 1,
		RequestID:     reqID,
		Type:          in.Type,
		ID:            in.ID,
		HalType:       in.HalType,
		Scenario:      in.Scenario,
	}

	// paramsなしケース
	if len(in.ParamsRaw) == 0 || bytes.Equal(bytes.TrimSpace(in.ParamsRaw), []byte("null")) {
		if in.RequireParams {
			return nil, errors.New("params is required but missing/null")
		}
		// Paramsはomitemptyで省略
		return marshalWithLimit(n, lim)
	}

	// paramsありケース：配列としてdecode & normalize
	var rawParams []ParamKV
	dec := json.NewDecoder(bytes.NewReader(in.ParamsRaw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&rawParams); err != nil {
		return nil, fmt.Errorf("params must be an array of {key,value}: %w", err)
	}

	if lim.MaxParams > 0 && len(rawParams) > lim.MaxParams {
		return nil, fmt.Errorf("too many params: %d > %d", len(rawParams), lim.MaxParams)
	}

	out := make([]ParamNV, 0, len(rawParams))
	index := map[string]int{} // for last-wins

	for i, p := range rawParams {
		key := strings.TrimSpace(p.Key)
		if key == "" {
			return nil, fmt.Errorf("params[%d].key is empty", i)
		}
		if lim.MaxKeyLen > 0 && utf8.RuneCountInString(key) > lim.MaxKeyLen {
			return nil, fmt.Errorf("params[%d].key too long", i)
		}

		valStr, err := normalizeValueToStringWithLimits(p.Value, lim)
		if err != nil {
			return nil, fmt.Errorf("params[%d].value invalid: %w", i, err)
		}
		valStr = strings.TrimSpace(valStr)

		if lim.MaxValueLen > 0 && utf8.RuneCountInString(valStr) > lim.MaxValueLen {
			return nil, fmt.Errorf("params[%d].value too long", i)
		}

		// duplicate policy
		if pos, ok := index[key]; ok {
			switch dup {
			case DupKeyError:
				return nil, fmt.Errorf("duplicate key: %s", key)
			case DupKeyLastWins:
				out[pos] = ParamNV{Key: key, Value: valStr}
			default:
				return nil, fmt.Errorf("unknown dup policy")
			}
		} else {
			index[key] = len(out)
			out = append(out, ParamNV{Key: key, Value: valStr})
		}
	}

	n.Params = out
	return marshalWithLimit(n, lim)
}

func marshalWithLimit(n Normalized, lim Limits) ([]byte, error) {
	b, err := json.Marshal(n)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal normalized json: %w", err)
	}
	if lim.MaxPayloadBytes > 0 && len(b) > lim.MaxPayloadBytes {
		return nil, fmt.Errorf("normalized payload too large: %d > %d", len(b), lim.MaxPayloadBytes)
	}
	return b, nil
}

//
// 7) value normalization (supports object like camera)
//

func normalizeValueToStringWithLimits(v any, lim Limits) (string, error) {
	switch x := v.(type) {
	case string:
		return x, nil
	case bool:
		if x {
			return "true", nil
		}
		return "false", nil
	case float64:
		// encoding/json は数値を float64 にする
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10), nil
		}
		return strconv.FormatFloat(x, 'g', -1, 64), nil
	case map[string]any:
		if err := validateContainer(x, 1, lim); err != nil {
			return "", err
		}
		b, err := json.Marshal(x)
		if err != nil {
			return "", err
		}
		return string(b), nil // ★cameraのようなobjectはJSON文字列に
	case []any:
		if err := validateContainer(x, 1, lim); err != nil {
			return "", err
		}
		b, err := json.Marshal(x)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case nil:
		return "", errors.New("null is not allowed")
	default:
		return "", fmt.Errorf("unsupported type %T", v)
	}
}

func validateContainer(v any, depth int, lim Limits) error {
	if lim.MaxDepth > 0 && depth > lim.MaxDepth {
		return fmt.Errorf("value nesting too deep: %d > %d", depth, lim.MaxDepth)
	}

	switch x := v.(type) {
	case map[string]any:
		if lim.MaxContainerLen > 0 && len(x) > lim.MaxContainerLen {
			return fmt.Errorf("object too large: %d > %d", len(x), lim.MaxContainerLen)
		}
		for _, vv := range x {
			switch vv.(type) {
			case map[string]any, []any:
				if err := validateContainer(vv, depth+1, lim); err != nil {
					return err
				}
			default:
				// scalarはOK
			}
		}
		return nil
	case []any:
		if lim.MaxContainerLen > 0 && len(x) > lim.MaxContainerLen {
			return fmt.Errorf("array too large: %d > %d", len(x), lim.MaxContainerLen)
		}
		for _, vv := range x {
			switch vv.(type) {
			case map[string]any, []any:
				if err := validateContainer(vv, depth+1, lim); err != nil {
					return err
				}
			default:
			}
		}
		return nil
	default:
		return nil
	}
}

//
// 8) cgo send
//

func sendToCpp(payload []byte) error {
	if len(payload) == 0 {
		return errors.New("normalized payload is empty")
	}
	rc := int(C.Process((*C.char)(unsafe.Pointer(&payload[0])), C.int(len(payload))))
	if rc != 0 {
		return fmt.Errorf("cpp Process failed: rc=%d", rc)
	}
	return nil
}


// 以下は使い方を記載

var ridGen = bridge.NewRequestIDGen(999999, 0)

var lim = bridge.Limits{
	MaxParams:       200,
	MaxKeyLen:       64,
	MaxValueLen:     4096,
	MaxPayloadBytes: 64 * 1024,
	MaxDepth:        5,
	MaxContainerLen: 200,
}

// 実行用
_ = bridge.NormalizeSendExecute(rawBody, lim, ridGen, bridge.DupKeyLastWins)

// 設定用
_ = bridge.NormalizeSendConfig(rawBody, lim, ridGen)
