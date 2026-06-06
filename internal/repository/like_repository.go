package repository

import (
	"context"
	"fmt"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
)

// いいね用のリポジトリ
type LikeRepository struct {
	db DBExecutor
}

// いいね用リポジトリのインスタンスを生成
func NewLikeRepository(db DBExecutor) *LikeRepository {
	return &LikeRepository{db: db}
}

// 投稿にいいねを追加する
func (r *LikeRepository) Create(ctx context.Context, userID int, postID int) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, postID)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to like post : PostID=%d", postID), err)
	}
	return nil
}

// 投稿のいいねを削除する
func (r *LikeRepository) Delete(ctx context.Context, userID int, postID int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to remove like : PostID=%d", postID), err)
	}
	return nil
}

// 指定した投稿のいいねユーザーID一覧を取得する
func (r *LikeRepository) ListUserIDsByPostID(ctx context.Context, postID int) ([]int, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id FROM likes WHERE post_id = $1", postID)
	if err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to fetch likes : PostID=%d", postID), err)
	}
	defer rows.Close()

	userIDs := []int{}
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to scan row : PostID=%d", postID), err)
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to fetch likes : PostID=%d", postID), err)
	}

	return userIDs, nil
}
