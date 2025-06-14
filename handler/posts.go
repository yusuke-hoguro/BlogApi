package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// Post構造体
type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int    `json:"user_id"`
}

// HTTPメソッドごとのルーティング
func PostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// 記事一覧取得
			GetPostsHandler(db)(w, r)
		case http.MethodPost:
			// 記事作成用
			CreatePostHandler(db)(w, r)
		case http.MethodPut:
			// 記事更新用
			UpdatePostHandler(db)(w, r)
		case http.MethodDelete:
			// 記事削除
			DeletePostHandler(db)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// 記事一覧取得用のハンドラー関数
func GetPostsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// postsテーブルから指定したカラムのデータを取得する
		rows, err := db.Query("SELECT id, title, content, user_id FROM posts")
		if err != nil {
			http.Error(w, "DBクエリエラー", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var p Post
			// DBから取得したデータをGoの変数に格納
			if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.UserID); err != nil {
				http.Error(w, "データ取得エラー", http.StatusInternalServerError)
				return
			}
			// Post構造体のスライスに追加
			posts = append(posts, p)
		}
		// json形式に変更してレスポンスに書き込む
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
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
		var post Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		var post Post
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
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

		w.WriteHeader(http.StatusNoContent)
		fmt.Fprintln(w, "Post deleted successfully!")
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
		rows, err := db.Query("SELECT id, title, content, user_id FROM posts WHERE user_id = $1", userID)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []Post
		for rows.Next() {
			var post Post
			if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
				http.Error(w, "Post not found", http.StatusNotFound)
				return
			}
			posts = append(posts, post)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
	}
}
