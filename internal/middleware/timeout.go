package middleware

import (
	"context"
	"net/http"
	"time"
)

// タイムアウトミドルウェア
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// タイムアウト付きのコンテキストを作成
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			// リクエストに新しいコンテキストを設定
			r = r.WithContext(ctx)
			// 次のハンドラーへ処理を渡す
			next.ServeHTTP(w, r)
		})
	}
}
