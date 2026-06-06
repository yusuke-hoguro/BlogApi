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

// コメントを作成する
func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	// コメントを挿入する
	query := `INSERT INTO comments (post_id, user_id, content) 
				VALUES ($1, $2, $3)
				RETURNING id, post_id, user_id, content, created_at`

	err := r.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert comment", err)
	}
	return nil
}

// 指定したコメントIDから所有者のユーザーIDと投稿IDを取得する
func (r *CommentRepository) FindOwnerByID(ctx context.Context, commentID int) (int, int, error) {
	var userID, postID int
	err := r.db.QueryRowContext(ctx, "SELECT user_id, post_id FROM comments WHERE id = $1", commentID).Scan(&userID, &postID)
	if err == sql.ErrNoRows {
		return 0, 0, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Comment not found : CommentID=%d", commentID), err)
	} else if err != nil {
		return 0, 0, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : CommentID=%d", commentID), err)
	}
	return userID, postID, nil
}

// 指定したIDのコメントを削除する
func (r *CommentRepository) Delete(ctx context.Context, commentID int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to delete comment : CommentID=%d", commentID), err)
	}
	return checkRowAffected(result, fmt.Sprintf("Comment not found : CommentID=%d", commentID))
}

// 指定したIDのコメントを更新する
func (r *CommentRepository) Update(ctx context.Context, commentID int, content string) error {
	result, err := r.db.ExecContext(ctx, "UPDATE comments SET content = $1 WHERE id = $2", content, commentID)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Failed to update comment : CommentID=%d", commentID), err)
	}
	return checkRowAffected(result, fmt.Sprintf("Comment not found : CommentID=%d", commentID))
}
