package service

import (
	"context"
	"fmt"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
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

// リクエストのユーザーとコメント所有者を確認する
func (s *CommentService) EnsureCommentOwner(ctx context.Context, userID int, commentID int) (int, error) {
	commentOwnerID, postID, err := s.repo.FindOwnerByID(ctx, commentID)
	if err != nil {
		return 0, err
	}
	if commentOwnerID != userID {
		return 0, apperror.NewAppError(apperror.TypeForbidden, fmt.Sprintf("Forbidden : CommentID=%d", commentID), nil)
	}
	return postID, nil
}

// コメントの削除処理を実施する
func (s *CommentService) DeleteComment(ctx context.Context, commentID int) error {
	return s.repo.Delete(ctx, commentID)
}

// コメントの更新処理を実施する
func (s *CommentService) UpdateComment(ctx context.Context, commentID int, content string) error {
	return s.repo.Update(ctx, commentID, content)
}
