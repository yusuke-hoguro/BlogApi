package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

// 【デバッグ用】タイムアウトミドルウェアの動作確認用ハンドラー　※routes.goには登録しないこと！
func SleepHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		select {
		// 5秒待機してからレスポンスを返す
		case <-time.After(30 * time.Second):
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("completed"))
			return
		// クライアントがリクエストをキャンセルした場合の処理
		case <-ctx.Done():
			err := ctx.Err()

			// タイムアウトの場合とキャンセルの場合でエラーメッセージを分ける
			if errors.Is(err, context.DeadlineExceeded) {
				http.Error(w, "Request timeout", http.StatusGatewayTimeout)
				log.Println("Request timed out")
				return
			}
			if errors.Is(err, context.Canceled) {
				http.Error(w, "Request cancelled", http.StatusRequestTimeout)
				log.Println("Request canceled by client")
				return
			}
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
			log.Printf("Request canceled: %v", err)
			return
		}
	})
}
