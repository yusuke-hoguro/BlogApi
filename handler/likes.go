package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// 指定した投稿に「いいね」をつける
func LikePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 認証情報からユーザーIDを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 「いいね」を登録する
		_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, postID)
		if err != nil {
			http.Error(w, "Failed to like post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Post liked successfully")
	}

}

// 指定した投稿の「いいね」数とユーザーを取得する
func GetLikesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 「いいね」の数とユーザー一覧を取得する
		rows, err := db.Query("SELECT user_id FROM likes WHERE post_id = $1", postID)
		if err != nil {
			http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		//「いいね」を押してくれたユーザーIDを保管する
		var userIDs []int
		for rows.Next() {
			var userID int
			if err := rows.Scan(&userID); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			userIDs = append(userIDs, userID)
		}

		// JSONレスポンスを返す
		resp := map[string]any{
			"post_id":    postID,
			"like_count": len(userIDs),
			"user_ids":   userIDs,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// 指定した投稿の「いいね」を削除する
func UnlikePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 認証情報からユーザーIDを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 「いいね」を削除する
		_, err = db.Exec("DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID)
		if err != nil {
			http.Error(w, "Failed to remove like", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "like removed successfully!")
		w.WriteHeader(http.StatusNoContent)
	}

}
