package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
)

// GetPostsByIDHandler godoc
// @Summary 投稿をIDで取得する
// @Description 指定したIDの投稿を返す
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Produce json
// @Param id path int true "PostID"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id} [get]
func GetPostsByIDHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// URLから投稿IDを抽出する
		id, appErr := postIDFromRequest(r)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// postsテーブルから指定したカラムのデータを取得する
		var post models.Post
		err := db.QueryRowContext(ctx, "SELECT id, title, content, user_id, created_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Post not found : PostID=%d", id), err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : PostID=%d", id), err))
			return
		}

		// 取得した投稿をJSONで返す
		respondJSON(w, http.StatusOK, post)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_fetched", UserID: post.UserID, PostID: post.ID})
	}
}

// CreatePostHandler godoc
// @Summary 新しい投稿を作成する
// @Description 送られてきた構造体のデータから新規投稿を作成する
// @Description
// @Description **エラー条件:**
// @Description - 無効な投稿内容、タイトルか投稿内容が空、タイトルが100文字以上、投稿内容が1000文字以上 → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param post body models.Post true "投稿内容"
// @Success 201 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts [post]
func CreatePostHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからユーザーIDを取得する
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// Post型の構造体にデコードして格納
		var post models.Post
		if appErr := decodeJSON(r, &post); appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// 投稿のバリデーションを行う
		if err := validatePostInput(post); err != nil {
			respondAppError(w, err)
			return
		}

		// 記事にユーザーIDを設定する
		post.UserID = userID

		// トランザクションを開始する
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to start transaction", err))
			return
		}
		defer func() {
			if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
				fmt.Printf("Failed to rollback transaction: %v\n", err)
			}
		}()

		// 投稿 INSERT実行
		err = tx.QueryRowContext(ctx, "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id", post.Title, post.Content, post.UserID).Scan(&post.ID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert post", err))
			return
		}

		// 投稿統計 INSERT実行
		_, err = tx.ExecContext(ctx, "INSERT INTO post_stats (post_id) VALUES ($1)", post.ID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert post stats", err))
			return
		}

		// トランザクションをコミットする
		if err := tx.Commit(); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to commit transaction", err))
			return
		}

		// 作成した投稿をJSONで返す
		respondJSON(w, http.StatusCreated, post)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_created", UserID: post.UserID, PostID: post.ID})
	}
}

// UpdatePostHandler godoc
// @Summary 投稿の内容を更新する
// @Description 送られてきた構造体のデータから投稿を更新する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID、無効な投稿内容、タイトルか投稿内容が空、タイトルが100文字以上、投稿内容が1000文字以上 → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 送信者が記事の投稿者でない → 403 Forbidden
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "投稿ID"
// @Param post body models.Post true "投稿内容"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id} [put]
func UpdatePostHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからユーザーIDを取得する
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// URLから投稿IDを抽出する
		id, appErr := postIDFromRequest(r)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// DBから投稿者のユーザーIDを取得する
		var postUserID int
		err := db.QueryRowContext(ctx, "SELECT user_id FROM posts WHERE id = $1", id).Scan(&postUserID)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Post not found : PostID=%d", id), err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : PostID=%d", id), err))
			return
		}

		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, fmt.Sprintf("Forbidden : PostID=%d", id), nil))
			return
		}

		// Post型の構造体にデコードして格納
		var post models.Post
		if appErr := decodeJSON(r, &post); appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// 投稿のバリデーションを行う
		if err := validatePostInput(post); err != nil {
			respondAppError(w, err)
			return
		}

		// UPDATE実行
		result, err := db.ExecContext(ctx, "UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to update post", err))
			return
		}

		// 更新行数の確認
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Post not found or no changes", err))
			return
		}

		// 更新した投稿をJSONで返す
		respondJSON(w, http.StatusOK, post)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_updated", UserID: userID, PostID: id})
	}
}

// DeletePostHandler godoc
// @Summary 投稿を削除する
// @Description 送られてきたIDの投稿を削除する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 送信者が記事の投稿者でない → 403 Forbidden
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "投稿ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id} [put]
func DeletePostHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからリクエストをなげたユーザーIDを取得
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// URLからIDを取得する
		id, appErr := postIDFromRequest(r)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// 削除対象の投稿を作成したユーザーのIDを取得する
		var postUserID int
		err := db.QueryRow("SELECT user_id FROM posts WHERE id = $1", id).Scan(&postUserID)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, fmt.Sprintf("Post not found : PostID=%d", id), err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, fmt.Sprintf("Database error : PostID=%d", id), err))
			return
		}

		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, fmt.Sprintf("Forbidden : PostID=%d", id), nil))
			return
		}

		// DELETE実行
		result, err := db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to delete post", err))
			return
		}

		// 削除行数の確認
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to confirm deletion", err))
			return
		} else if rowsAffected == 0 {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Post not found", nil))
			return
		}

		// 削除成功のため204 No Contentを返す
		respondJSON(w, http.StatusNoContent, nil)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_deleted", UserID: userID, PostID: id})
	}
}

// GetMyPostsHandler godoc
// @Summary ユーザー自身の投稿を取得する
// @Description リクエストを投げたユーザーが作成した投稿を取得する
// @Description
// @Description **エラー条件:**
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {array} models.Post
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/myposts [get]
func GetMyPostsHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからリクエストをなげたユーザーIDを取得
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// DBから自分の投稿を取得
		rows, err := db.QueryContext(ctx, "SELECT id, title, content, user_id, created_at FROM posts WHERE user_id = $1", userID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch posts", err))
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
				respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to parse post", err))
				return
			}
			posts = append(posts, post)
		}

		// rows.Next()のループが終了した後にエラーが発生していないか確認する(DBからのデータ取得中にエラーが発生していないか)
		if err := rows.Err(); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch posts", err))
			return
		}

		// 取得した投稿をJSONで返す
		respondJSON(w, http.StatusOK, posts)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "my_posts_fetched", UserID: userID})
	}
}

// GetAllPostsHandler godoc
// @Summary すべての投稿を取得する
// @Description DBから全投稿を取得して返却する
// @Description 送られてきたIDの投稿を削除する
// @Description
// @Description **エラー条件:**
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts [get]
func GetAllPostsHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// 全投稿を取得する
		rows, err := db.QueryContext(ctx, "SELECT id, title, content, user_id, created_at FROM posts ORDER BY created_at DESC")
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch posts", err))
			return
		}
		defer rows.Close()

		// nilをJSON化しないようにスライスを初期化する
		posts := []models.Post{}

		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
				respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to parse post", err))
				return
			}
			posts = append(posts, post)
		}

		// rows.Next()のループが終了した後にエラーが発生していないか確認する(DBからのデータ取得中にエラーが発生していないか)
		if err := rows.Err(); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch posts", err))
			return
		}

		// 取得した投稿をJSONで返す
		respondJSON(w, http.StatusOK, posts)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "posts_fetched"})
	}
}
