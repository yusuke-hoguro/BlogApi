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
	Like    *service.LikeService
	User    *service.UserService
}

// サービスの初期化を行う関数
func NewServices(db *sql.DB) *Services {
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	userRepo := repository.NewUserRepository(db)

	return &Services{
		Post:    service.NewPostService(postRepo),
		Comment: service.NewCommentService(commentRepo),
		Like:    service.NewLikeService(likeRepo),
		User:    service.NewUserService(userRepo),
	}
}
