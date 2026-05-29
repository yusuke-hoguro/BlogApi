package repository

import (
	"database/sql"
)

// 投稿用のリポジトリ
type CommentRepository struct {
	db *sql.DB
}

// 投稿用リポジトリのインスタンスを生成
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}
