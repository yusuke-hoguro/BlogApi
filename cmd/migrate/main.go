package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
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

	// マイグレーション実行時のタイムアウト時間を取得
	timeout, err := migrationTimeout()
	if err != nil {
		return err
	}

	// 設定されたタイムアウトで自動キャンセルするコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// DBマイグレーションを実行する
	if err := db.RunMigrations(ctx, conn, "sql/migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("migration completed")
	return nil
}

// DBマイグレーションのタイムアウト時間を取得する
func migrationTimeout() (time.Duration, error) {
	const defaultTimeoutSeconds = 300
	// 環境変数からマイグレーション実行のタイムアウト時間を取得する
	value := os.Getenv("MIGRATION_TIMEOUT_SECONDS")
	if value == "" {
		return time.Duration(defaultTimeoutSeconds) * time.Second, nil
	}
	// 取得した文字列を秒数に変換する
	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid MIGRATION_TIMEOUT_SECONDS: %w", err)
	}
	if seconds <= 0 {
		return 0, fmt.Errorf("MIGRATION_TIMEOUT_SECONDS must be positive")
	}
	return time.Duration(seconds) * time.Second, nil
}
