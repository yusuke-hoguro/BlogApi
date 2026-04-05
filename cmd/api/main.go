// @title Blog API
// @version 1.0
// @description This is a sample blog API built with Go net/http.
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"context"
	"errors"
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
	"golang.org/x/sync/errgroup"

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
	// ルートの登録
	router.RegisterRoutes(r, conn)
	// CORSミドルウェアを適用
	handler := middleware.CorsMiddleware(r)
	// タイムアウトミドルウェアを適用(戻り値が関数なので（handler）をつけて実行する)
	handler = middleware.TimeoutMiddleware(10 * time.Second)(handler)

	// HTTPサーバーの設定
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// シグナルを受け取るためのコンテキストを作成
	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// errgroupでgoroutineのエラー管理とキャンセル伝播を行う
	g, ctx := errgroup.WithContext(sigCtx)

	// サーバー起動を起動するgoroutine
	g.Go(func() error {
		return runHTTPServer(srv)
	})

	// コンテキストがキャンセルされたらサーバーをシャットダウンするgoroutine
	g.Go(func() error {
		return shutdownOnContextDone(ctx, srv)
	})

	// いずれかのgoroutineがエラーを返すのを待つ
	if err := g.Wait(); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// HTTPサーバーを起動する
func runHTTPServer(srv *http.Server) error {
	log.Printf("Server started at %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// コンテキストがキャンセルされたらサーバーをシャットダウンする
func shutdownOnContextDone(ctx context.Context, srv *http.Server) error {
	<-ctx.Done()
	log.Printf("Shutdown signal received: %v", ctx.Err())

	// shutdownのタイムアウトを設定
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// サーバーをシャットダウン
	log.Printf("Server shutting down...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("Server shutdown failed: %w", err)
	}
	log.Println("Server shutdown complete")
	return nil
}
