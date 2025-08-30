package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/yusuke-hoguro/BlogApi/middleware"
	"github.com/yusuke-hoguro/BlogApi/models"
)

const (
	maxTitleLength   = 100
	maxContentLength = 1000
)

// 記事一覧取得用のハンドラー関数
func GetPostsByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// IDを抽出する
		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		// postsテーブルから指定したカラムのデータを取得する
		var post models.Post
		err = db.QueryRow("SELECT id, title, content, user_id, created_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// json形式に変更してレスポンスに書き込む
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(post)
	}
}

// 記事作成用のハンドラー関数
func CreatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからユーザーIDを取得する
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Post型の構造体にデコードして格納
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// タイトルが空の場合はエラーとする
		if strings.TrimSpace(post.Title) == "" {
			http.Error(w, "Title must not be empty", http.StatusBadRequest)
			return
		}

		// タイトルが100文字より大きい場合はエラーとする
		if len(post.Title) > maxTitleLength {
			http.Error(w, "Title must be 100 characters or less", http.StatusBadRequest)
			return
		}

		// 投稿の内容が空の場合はエラーとする
		if strings.TrimSpace(post.Content) == "" {
			http.Error(w, "Content is required", http.StatusBadRequest)
			return
		}

		// タイトルが1000文字より大きい場合はエラーとする
		if len(post.Content) > maxContentLength {
			http.Error(w, "Content must be 1000 characters or less", http.StatusBadRequest)
			return
		}

		// 記事にユーザーIDを設定する
		post.UserID = userID

		// INSERT実行
		err := db.QueryRow("INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3) RETURNING id", post.Title, post.Content, post.UserID).Scan(&post.ID)
		if err != nil {
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		// 作成した記事IDをJSONで返す
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	}
}

// 記事更新のハンドラー関数
func UpdatePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからユーザーIDを取得する
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// URLからIDを取得する
		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// DBから投稿者のユーザーIDを取得する
		var postUserID int
		err = db.QueryRow("SELECT user_id FROM posts WHERE id = $1", id).Scan(&postUserID)
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Post型の構造体にデコードして格納
		var post models.Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// タイトルが空の場合はエラーとする
		if strings.TrimSpace(post.Title) == "" {
			http.Error(w, "Title must not be empty", http.StatusBadRequest)
			return
		}

		// タイトルが100文字より大きい場合はエラーとする
		if len(post.Title) > maxTitleLength {
			http.Error(w, "Title must be 100 characters or less", http.StatusBadRequest)
			return
		}

		// 投稿の内容が空の場合はエラーとする
		if strings.TrimSpace(post.Content) == "" {
			http.Error(w, "Content is required", http.StatusBadRequest)
			return
		}

		// タイトルが1000文字より大きい場合はエラーとする
		if len(post.Content) > maxContentLength {
			http.Error(w, "Content must be 1000 characters or less", http.StatusBadRequest)
			return
		}

		// UPDATE実行
		result, err := db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", post.Title, post.Content, id)
		if err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		// 更新行数の確認
		rowsAffected, err := result.RowsAffected()
		if err != nil || rowsAffected == 0 {
			http.Error(w, "Post nor found or no changes", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Post update successfully!")
	}
}

// 記事削除のハンドラー関数
func DeletePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからリクエストをなげたユーザーIDを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// URLからIDを取得する
		idStr := strings.TrimPrefix(r.URL.Path, "/posts/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 削除対象の投稿を作成したユーザーのIDを取得する
		var postUserID int
		err = db.QueryRow("SELECT user_id FROM posts WHERE id = $1", id).Scan(&postUserID)
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// DELETE実行
		result, err := db.Exec("DELETE FROM posts WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Failed to update post", http.StatusInternalServerError)
			return
		}

		// 削除行数の確認
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Failed to confirm deletion", http.StatusInternalServerError)
			return
		} else if rowsAffected == 0 {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		fmt.Fprintln(w, "Post deleted successfully!")
		w.WriteHeader(http.StatusNoContent)
	}
}

// ユーザー自身の投稿のみを取得する
func GetMyPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからリクエストをなげたユーザーIDを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// DBから自分の投稿を取得
		rows, err := db.Query("SELECT id, title, content, user_id, created_at FROM posts WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}

// 全投稿を取得する
func GetAllPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 全投稿を取得する
		rows, err := db.Query("SELECT id, title, content, user_id, created_at FROM posts ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var post models.Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt); err != nil {
				http.Error(w, "Failed to parse post", http.StatusNotFound)
				return
			}
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
