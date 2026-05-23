package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yusuke-hoguro/BlogApi/internal/db"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Printf("migration command failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// .envファイルを読み込む
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env: %v", err)
	}

	// DBに接続する
	conn, err := db.ConnectDB()
	if err != nil {
		return fmt.Errorf("DB接続に失敗: %w", err)
	}
	defer conn.Close()

	// 30秒で自動キャンセルするコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// DBマイグレーションを実行する
	if err := db.RunMigrations(ctx, conn, "sql/migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("migration completed")
	return nil
}
