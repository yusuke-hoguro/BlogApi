package app

import (
	"database/sql"

	"github.com/yusuke-hoguro/BlogApi/internal/repository"
	"github.com/yusuke-hoguro/BlogApi/internal/service"
)

// サービスをまとめる構造体
type Services struct {
	Post    *service.PostService
	Comment *service.CommentService
}

// サービスの初期化を行う関数
func NewServices(db *sql.DB) *Services {
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	return &Services{
		Post:    service.NewPostService(postRepo),
		Comment: service.NewCommentService(commentRepo),
	}
}
