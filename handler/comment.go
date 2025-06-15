package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/middleware"
	"github.com/yusuke-hoguro/BlogApi/models"
)

// 投稿のコメント取得用ハンドラー関数
func GetCommentsByPostIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 指定した投稿のコメントを取得する
		rows, err := db.Query(`
			SELECT id, post_id, user_id, content, created_at
			FROM comments
			WHERE post_id = $1
			ORDER BY created_at ASC
		`, postID)
		if err != nil {
			http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 取得したコメントをスライスに格納
		var comments []models.Comment
		for rows.Next() {
			var c models.Comment
			if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
				http.Error(w, "Error reading comment", http.StatusInternalServerError)
				return
			}
			comments = append(comments, c)
		}

		// コメントがない場合
		if len(comments) == 0 {
			fmt.Println("No comments.")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	}
}

// コメント投稿用のハンドラー関数
func PostCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからuser_idを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// リクエストボディからコメントを読み取る
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// コメントを挿入する
		query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)`
		_, err = db.Exec(query, postID, userID, comment.Content)
		if err != nil {
			http.Error(w, "Failed to insert comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Comment created"})

	}
}
