package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/models"
)

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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	}
}
