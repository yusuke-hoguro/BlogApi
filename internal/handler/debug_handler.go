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
			if _, err := w.Write([]byte("completed")); err != nil {
				log.Println("write error:", err)
			}
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

// 【デバッグ用】ゴルーチン待ちハンドラー関数 ※routes.goには登録しないこと！
func GorutineWaitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		resultCh := make(chan string, 1)

		// 5秒後に結果を返すゴルーチン
		go func() {
			select {
			case <-time.After(5 * time.Second):
				select {
				case resultCh <- "Background task completed":
				case <-ctx.Done():
					log.Println("Background task canceled before sending result")
					return
				}
			case <-ctx.Done():
				log.Println("Background task canceled")
				return
			}
		}()

		select {
		// 結果をチャネルから受け取る
		case result := <-resultCh:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(result))
			return
		// クライアントがリクエストをキャンセルした場合の処理
		case <-ctx.Done():
			err := ctx.Err()
			// タイムアウトの場合
			if errors.Is(err, context.DeadlineExceeded) {
				// サーバー側のタイムアウトなので504 Gateway Timeoutを返す
				http.Error(w, "Request timeout", http.StatusGatewayTimeout)
				log.Println("Request timed out")
				return
			}
			// クライアント側のキャンセルの場合
			if errors.Is(err, context.Canceled) {
				http.Error(w, "Request cancelled", http.StatusRequestTimeout)
				log.Println("Request canceled by client")
				return
			}
			// その他のエラーの場合もクライアントキャンセルとみなす
			http.Error(w, "Request cancelled", http.StatusRequestTimeout)
			log.Printf("Request canceled: %v", err)
			return
		}
	}
}
