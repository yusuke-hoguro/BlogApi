package handler

import (
	"log"
	"net/http"
)

// HealthzHandler godoc
// @Summary ヘルスチェック
// @Description BlogAPI が HTTP リクエストを受け付けられる状態か確認します
// @Description
// @Description **エラー条件:**
// @Description - 無効なメソッドタイプ → 405 Method Not Allowed
// @Tags health
// @Produce plain
// @Success 200 {string} string "OK"
// @Failure 405 {object} models.ErrorResponse
// @Router /api/healthz [get]
func HealthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			respondError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Printf("failed to write response: %v", err)
			return
		}
	}
}
