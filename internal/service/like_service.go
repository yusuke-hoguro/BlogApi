package service

import (
	"context"

	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/repository"
)

// いいね用サービスの構造体
type LikeService struct {
	repo *repository.LikeRepository
}

// いいね用サービスのインスタンスを生成する関数
func NewLikeService(repo *repository.LikeRepository) *LikeService {
	return &LikeService{repo: repo}
}

// 投稿にいいねを追加する
func (s *LikeService) LikePost(ctx context.Context, userID int, postID int) error {
	return s.repo.Create(ctx, userID, postID)
}

// 投稿のいいねを削除する
func (s *LikeService) UnlikePost(ctx context.Context, userID int, postID int) error {
	return s.repo.Delete(ctx, userID, postID)
}

// 投稿のいいね情報を取得する
func (s *LikeService) GetLikes(ctx context.Context, postID int) (*models.LikesResponse, error) {
	userIDs, err := s.repo.ListUserIDsByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return &models.LikesResponse{
		PostID:    postID,
		LikeCount: len(userIDs),
		UserIDs:   userIDs,
	}, nil
}
