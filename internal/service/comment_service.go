package service

import (
	"context"

	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/repository"
)

// コメント用サービスの構造体
type CommentService struct {
	repo *repository.CommentRepository
}

// コメント用サービスのインスタンスを生成する関数
func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

// 指定した投稿IDのコメントを取得する
func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
	return s.repo.ListByPostID(ctx, postID)
}

// 指定したIDのコメントを取得する
func (s *CommentService) GetCommentByID(ctx context.Context, commentID int) (*models.Comment, error) {
	return s.repo.FindByID(ctx, commentID)
}

// コメントの作成処理を実施する
func (s *CommentService) CreateComment(ctx context.Context, postID int, userID int, comment *models.Comment) error {
	comment.PostID = postID
	comment.UserID = userID
	return s.repo.Create(ctx, comment)
}
