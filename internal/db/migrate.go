package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 実行するSQL文を定数で定義する
const createSchemaMigrationsTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version TEXT PRIMARY KEY,
	applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

// すべてのマイグレーションを実行する関数
func RunMigrations(ctx context.Context, conn *sql.DB, migrationsDir string) error {
	// マイグレーション管理用のテーブルを作成する
	if _, err := conn.ExecContext(ctx, createSchemaMigrationsTable); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// マイグレーションファイルのリストを取得する
	files, err := migrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	// マイグレーションファイルを順番に実行する
	for _, file := range files {
		if err := runMigration(ctx, conn, file); err != nil {
			return err
		}
	}
	return nil
}

// マイグレーションファイルを実行する関数
func runMigration(ctx context.Context, conn *sql.DB, file string) error {
	// マイグレーションファイルのバージョンを取得する
	version := filepath.Base(file)

	// マイグレーションが既に適用されているかどうかを確認する
	applied, err := isMigrationApplied(ctx, conn, version)
	if err != nil {
		return err
	}

	if applied {
		fmt.Printf("skip migration: %s\n", version)
		return nil
	}

	sqlText, err := readMigrationSQL(file)
	if err != nil {
		return err
	}

	if err := applyMigration(ctx, conn, version, sqlText); err != nil {
		return err
	}

	fmt.Printf("applied migration: %s\n", version)
	return nil
}

// マイグレーションファイルの内容を読み込む関数
func readMigrationSQL(file string) (string, error) {
	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read migration file %s: %w", file, err)
	}

	sqlText := strings.TrimSpace(string(sqlBytes))
	if sqlText == "" {
		return "", fmt.Errorf("migration %s is empty", filepath.Base(file))
	}

	return sqlText, nil
}

// マイグレーションSQLを実行して適用済みとして記録する関数
func applyMigration(ctx context.Context, conn *sql.DB, version string, sqlText string) error {
	// トランザクションを開始する
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start migration transaction %s: %w", version, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			fmt.Printf("failed to rollback migration %s: %v\n", version, err)
		}
	}()

	// マイグレーションを実行する
	if _, err := tx.ExecContext(ctx, sqlText); err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", version, err)
	}

	// マイグレーションのバージョンを記録する
	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}

	// トランザクションをコミットする
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", version, err)
	}

	return nil
}

// マイグレーションファイルのリストを取得する関数
func migrationFiles(migrationsDir string) ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		// ディレクトリやSQL以外のファイルはスキップする
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		// マイグレーションファイルのパスをリストに追加する
		files = append(files, filepath.Join(migrationsDir, entry.Name()))
	}
	// ファイル名でソートする（マイグレーションの順序を保証するため）
	sort.Strings(files)
	return files, nil
}

// マイグレーションが既に適用されているかどうかを確認する関数
func isMigrationApplied(ctx context.Context, conn *sql.DB, version string) (bool, error) {
	var exists bool
	err := conn.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check migration %s: %w", version, err)
	}
	return exists, nil
}
