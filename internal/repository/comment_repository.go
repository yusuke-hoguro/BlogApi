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

// 指定したIDのコメントを見つける(存在しない場合はnilを返したいのでポインタを返す)
func (r *CommentRepository) FindByID(ctx context.Context, id int) (*models.Comment, error) {
	// 指定したIDのコメントを取得する
	var comment models.Comment
	err := r.db.QueryRowContext(ctx, "SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1", id).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Comment Not Found : CommentID=%d", id), err)
	} else if err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : CommentID=%d", id), err)
	}
	return &comment, nil
}
