package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
)

const (
	MaxTitleLength   = 100
	MaxContentLength = 1000
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
func GetPostsByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// URLから投稿IDを抽出する
		idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid ID : PostID="+idStr, err))
			return
		}

		// postsテーブルから指定したカラムのデータを取得する
		var post models.Post
		err = db.QueryRowContext(ctx, "SELECT id, title, content, user_id, created_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Post not found : PostID="+idStr, err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Database error : PostID="+idStr, err))
			return
		}

		// json形式に変更してレスポンスに書き込む
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(post); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to encode response", err))
			return
		}
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
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからユーザーIDを取得する
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized userID not found in context", nil))
			return
		}

		// Post型の構造体にデコードして格納
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body", err))
			return
		}

		// タイトルが空の場合はエラーとする
		if strings.TrimSpace(post.Title) == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Title must not be empty", nil))
			return
		}

		// タイトルが100文字より大きい場合はエラーとする
		if len(post.Title) > MaxTitleLength {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Title must be 100 characters or less", nil))
			return
		}

		// 投稿の内容が空の場合はエラーとする
		if strings.TrimSpace(post.Content) == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content is required", nil))
			return
		}

		// 投稿内容が1000文字より大きい場合はエラーとする
		if len(post.Content) > MaxContentLength {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content must be 1000 characters or less", nil))
			return
		}

		// 記事にユーザーIDを設定する
		post.UserID = userID

		// INSERT実行
		err := db.QueryRowContext(ctx, "INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id", post.Title, post.Content, post.UserID).Scan(&post.ID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert post", err))
			return
		}

		// 作成した記事IDをJSONで返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(post); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to encode response", err))
			return
		}
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
func UpdatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからユーザーIDを取得する
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized userID not found in context", nil))
			return
		}

		// URLからIDを取得する
		idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid ID : PostID="+idStr, err))
			return
		}

		// DBから投稿者のユーザーIDを取得する
		var postUserID int
		err = db.QueryRowContext(ctx, "SELECT user_id FROM posts WHERE id = $1", id).Scan(&postUserID)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Post not found : PostID="+idStr, err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Database error : PostID="+idStr, err))
			return
		}

		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, "Forbidden : PostID="+idStr, nil))
			return
		}

		// Post型の構造体にデコードして格納
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body", err))
			return
		}

		// タイトルが空の場合はエラーとする
		if strings.TrimSpace(post.Title) == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Title must not be empty", nil))
			return
		}

		// タイトルが100文字より大きい場合はエラーとする
		if len(post.Title) > MaxTitleLength {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Title must be 100 characters or less", nil))
			return
		}

		// 投稿の内容が空の場合はエラーとする
		if strings.TrimSpace(post.Content) == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content is required", nil))
			return
		}

		// タイトルが1000文字より大きい場合はエラーとする
		if len(post.Content) > MaxContentLength {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content must be 1000 characters or less", nil))
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(post); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからリクエストをなげたユーザーIDを取得
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized userID not found in context", nil))
			return
		}

		// URLからIDを取得する
		idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid ID : PostID="+idStr, err))
			return
		}

		// 削除対象の投稿を作成したユーザーのIDを取得する
		var postUserID int
		err = db.QueryRow("SELECT user_id FROM posts WHERE id = $1", id).Scan(&postUserID)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Post not found", err))
			return
		} else if err != nil {
			respondError(w, "Database error : PostID="+idStr, http.StatusInternalServerError)
			return
		}

		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, "Forbidden : PostID="+idStr, nil))
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

		if _, err := fmt.Fprintln(w, "Post deleted successfully!"); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
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
func GetMyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからリクエストをなげたユーザーIDを取得
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized userID not found in context", nil))
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
			return
		}
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
func GetAllPostsHandler(db *sql.DB) http.HandlerFunc {
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(posts); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
			return
		}
	}
}
