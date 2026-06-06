package repository

import (
	"context"
	"database/sql"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
)

// ユーザー用のリポジトリ
type UserRepository struct {
	db DBExecutor
}

// ユーザー用リポジトリのインスタンスを生成
func NewUserRepository(db DBExecutor) *UserRepository {
	return &UserRepository{db: db}
}

// ユーザーを作成する
func (r *UserRepository) Create(ctx context.Context, username string, hashedPassword string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, hashedPassword).Scan(&id)
	if err != nil {
		return 0, apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert user : Username="+username, err)
	}
	return id, nil
}

// ユーザー名から認証情報を取得する
func (r *UserRepository) FindAuthByUsername(ctx context.Context, username string) (int, string, error) {
	var id int
	var hashedPassword string
	err := r.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE username = $1", username).Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, "", apperror.NewAppError(apperror.TypeUnauthorized, "Invalid username or password : Username="+username, err)
	} else if err != nil {
		return 0, "", apperror.NewAppError(apperror.TypeInternalServer, "Database error : Username="+username, err)
	}
	return id, hashedPassword, nil
}
