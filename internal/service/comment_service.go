package service

import (
	"github.com/yusuke-hoguro/BlogApi/internal/repository"
)

// 投稿用サービスの構造体
type CommentService struct {
	repo *repository.CommentRepository
}

// 投稿用サービスのインスタンスを生成する関数
func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}
