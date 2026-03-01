// @title Blog API
// @version 1.0
// @description This is a sample blog API built with Go net/http.
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/internal/db"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/internal/router"

	_ "github.com/lib/pq"
)

func main() {
	if err := runServer(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func runServer() error {
	// DB接続を実施
	conn, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("DB接続失敗: %w", err)
	}
	defer conn.Close()
	// ポート取得
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ルーターの設定
	r := mux.NewRouter()
	router.RegisterRoutes(r, conn)
	// CORSミドルウェアを適用
	handler := middleware.CorsMiddleware(r)

	// SIGINT/SIGTERMを受けたら停止するコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// HTTPサーバーの設定
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// エラー用のチャネル
	errChan := make(chan error, 1)

	// サーバーを別ゴルーチンで起動
	go func() {
		log.Println("Server started at " + srv.Addr)
		errChan <- srv.ListenAndServe()
	}()

	// サーバー停止のシグナルを待つ
	select {
	case err := <-errChan: // サーバーエラーが発生した場合
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	case <-ctx.Done(): // SIGINT/SIGTERMを受けた場合
		log.Printf("Shutdown signal received: %v", ctx.Err())
	}

	// shutdownのタイムアウトを設定
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// サーバーをシャットダウン
	log.Printf("Server shutting down...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("Server shutdown failed: %w", err)
	}

	// サーバーが完全に停止するのを待つ
	if err := <-errChan; err != nil && err != http.ErrServerClosed {
		return err
	}

	log.Println("Server shutdown complete")
	return nil
}
