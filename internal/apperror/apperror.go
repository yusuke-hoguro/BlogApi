package apperror

import "net/http"

type Type string

// エラーの種類を定義
const (
	TypeBadRequest       Type = "bad_request"
	TypeUnauthorized     Type = "unauthorized"
	TypeForbidden        Type = "forbidden"
	TypeNotFound         Type = "not_found"
	TypeConflict         Type = "conflict"
	TypeTimeout          Type = "timeout"
	TypeInternalServer   Type = "internal_server_error"
	TypeMethodNotAllowed Type = "method_not_allowed"
)

// エラー構造体
type AppError struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
	Err     error  `json:"err,omitempty"`
}

// エラーインターフェースを実装
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return string(e.Type)
}

// 元々のエラー原因を返す
func (e *AppError) Unwrap() error {
	return e.Err
}

// エラーの生成関数
func NewAppError(errType Type, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// エラーのHTTPステータスコードを取得する関数
func GetStatusCode(errType Type) int {
	switch errType {
	case TypeBadRequest:
		return http.StatusBadRequest
	case TypeUnauthorized:
		return http.StatusUnauthorized
	case TypeForbidden:
		return http.StatusForbidden
	case TypeNotFound:
		return http.StatusNotFound
	case TypeConflict:
		return http.StatusConflict
	case TypeTimeout:
		return http.StatusRequestTimeout
	case TypeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}
