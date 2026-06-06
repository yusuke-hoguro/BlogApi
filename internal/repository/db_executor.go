package repository

import (
	"context"
	"database/sql"
)

// DBExecutor は sql.DB / sql.Tx の共通DB操作を表す。
type DBExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
