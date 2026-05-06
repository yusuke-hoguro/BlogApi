package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
)

// 投稿用のリポジトリ
type PostRepository struct {
	db *sql.DB
}

// 投稿用リポジトリのインスタンスを生成
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// 指定したIDから投稿を見つける(存在しない場合はnilを返したいのでポインタを返す)
func (r *PostRepository) FindByID(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	err := r.db.QueryRowContext(ctx, "SELECT id, title, content, user_id, created_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Post not found : PostID=%d", id), err)
	} else if err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : PostID=%d", id), err)
	}
	return &post, nil
}

// 指定したUserIDから投稿を見つける
func (r *PostRepository) ListByUserID(ctx context.Context, userID int) ([]models.Post, error) {
	return r.listPosts(ctx, "SELECT id, title, content, user_id, created_at FROM posts WHERE user_id = $1", userID)
}

// 全ての投稿を見つける
func (r *PostRepository) ListAll(ctx context.Context) ([]models.Post, error) {
	return r.listPosts(ctx, "SELECT id, title, content, user_id, created_at FROM posts ORDER BY created_at DESC")
}

// 指定したクエリを実行して投稿を見つける
func (r *PostRepository) listPosts(ctx context.Context, query string, args ...any) ([]models.Post, error) {
	// 指定されたクエリを実行する
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch posts", err)
	}
	defer rows.Close()

	// nilをJSON化しないようにスライスを初期化する
	posts := []models.Post{}
	// クエリの結果をスキャンしてpostsスライスに追加する
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
			return nil, apperror.NewAppError(apperror.TypeInternalServer, "Failed to parse post", err)
		}
		posts = append(posts, post)
	}

	// rows.Next()のループが終了した後にエラーが発生していないか確認する(DBからのデータ取得中にエラーが発生していないか)
	if err := rows.Err(); err != nil {
		return nil, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch posts", err)
	}
	return posts, nil
}

// 新しい投稿を作成する
func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	// トランザクションを開始する
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to start transaction", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			fmt.Printf("Failed to rollback transaction: %v\n", err)
		}
	}()
	// 投稿 INSERT実行
	err = tx.QueryRowContext(ctx, "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id", post.Title, post.Content, post.UserID).Scan(&post.ID)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert post", err)
	}
	// 投稿統計 INSERT実行
	_, err = tx.ExecContext(ctx, "INSERT INTO post_stats (post_id) VALUES ($1)", post.ID)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert post stats", err)
	}
	// トランザクションをコミットする
	if err := tx.Commit(); err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to commit transaction", err)
	}
	return nil
}

// 指定した投稿のIDからユーザーIDを見つける
func (r *PostRepository) FindUserIDByPostID(ctx context.Context, postID int) (int, error) {
	var userID int
	err := r.db.QueryRowContext(ctx, "SELECT user_id FROM posts WHERE id = $1", postID).Scan(&userID)
	if err == sql.ErrNoRows {
		return 0, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Post not found : PostID=%d", postID), err)
	} else if err != nil {
		return 0, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : PostID=%d", postID), err)
	}
	return userID, nil
}

// 指定したIDの投稿を更新する
func (r *PostRepository) Update(ctx context.Context, id int, post *models.Post) error {
	// UPDATE実行
	result, err := r.db.ExecContext(ctx, "UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to update post", err)
	}
	// 更新行数の確認
	return checkRowAffected(result, fmt.Sprintf("Post not found : PostID=%d", id))
}

// 指定したIDの投稿を削除する
func (r *PostRepository) Delete(ctx context.Context, id int) error {
	// DELETE実行
	result, err := r.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to delete post", err)
	}
	// 削除行数の確認
	return checkRowAffected(result, fmt.Sprintf("Post not found : PostID=%d", id))
}

// SQLの実行結果から影響を受けた行数を確認する関数
func checkRowAffected(result sql.Result, notFoundMessage string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return apperror.NewAppError(apperror.TypeInternalServer, "Failed to confirm operation", err)
	} else if rowsAffected == 0 {
		return apperror.NewAppError(apperror.TypeNotFound, notFoundMessage, nil)
	}
	return nil
}
