package service

import (
	"context"
	"fmt"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/repository"
)

// 投稿用サービスの構造体
type PostService struct {
	repo *repository.PostRepository
}

// 投稿用サービスのインスタンスを生成する関数
func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// リクエストのユーザーと投稿者を確認する
func (s *PostService) EnsurePostOwner(ctx context.Context, userID int, postID int) error {
	// DBから投稿者のユーザーIDを取得する
	postUserID, err := s.repo.FindUserIDByPostID(ctx, postID)
	if err != nil {
		return err
	}
	// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
	if postUserID != userID {
		return apperror.NewAppError(apperror.TypeForbidden, fmt.Sprintf("Forbidden : PostID=%d", postID), nil)
	}
	return nil
}

// 投稿の更新処理を実施する
func (s *PostService) UpdatePost(ctx context.Context, postID int, post *models.Post) error {
	return s.repo.Update(ctx, postID, post)
}

// 投稿の削除処理を実施する
func (s *PostService) DeletePost(ctx context.Context, postID int) error {
	return s.repo.Delete(ctx, postID)
}

// 投稿の作成処理を実施する
func (s *PostService) CreatePost(ctx context.Context, post *models.Post) error {
	return s.repo.Create(ctx, post)
}

// 指定した投稿IDの投稿を取得する
func (s *PostService) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {
	return s.repo.FindByID(ctx, postID)
}

// 指定したユーザーIDの投稿を全て取得する
func (s *PostService) GetPostsByUserID(ctx context.Context, userID int) ([]models.Post, error) {
	return s.repo.ListByUserID(ctx, userID)
}

// 全ての投稿を取得する
func (s *PostService) GetAllPosts(ctx context.Context) ([]models.Post, error) {
	return s.repo.ListAll(ctx)
}
