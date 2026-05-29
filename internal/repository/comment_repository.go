package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
)

// コメント用のリポジトリ
type CommentRepository struct {
	db *sql.DB
}

// コメント用リポジトリのインスタンスを生成
func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// 指定した投稿IDのコメントを見つける
func (r *CommentRepository) ListByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
	// 投稿IDを指定してコメントを取得する
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, post_id, user_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to fetch comments : PostID=%d", postID), err)
	}
	defer rows.Close()

	// 取得したコメントをスライスに格納(nilをJSON化しないようにスライスを初期化)
	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
			return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Error reading comment : PostID=%d", postID), err)
		}
		comments = append(comments, c)
	}

	// rows.Next()のループが終了した後にエラーが発生していないか確認する(DBからのデータ取得中にエラーが発生していないか)
	if err := rows.Err(); err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to fetch comments : PostID=%d", postID), err)
	}

	return comments, nil
}
