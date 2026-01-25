package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/models"
)

// LikePostHandler godoc
// @Summary 投稿に「いいね」をつける
// @Description 指定したIDの投稿に「いいね」をつける
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Produce json
// @Param id path int true "投稿ID"
// @Success 201 {string} string "Post liked successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/like [post]
func LikePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 認証情報からユーザーIDを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized", http.StatusUnauthorized)
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 「いいね」を登録する
		_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", userID, postID)
		if err != nil {
			respondError(w, "Failed to like post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if _, err := fmt.Fprintln(w, "Post liked successfully"); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetLikesHandler godoc
// @Summary 投稿の「いいね」を取得する
// @Description 指定したIDの投稿についている「いいね」を取得する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Produce json
// @Param id path int true "投稿ID"
// @Success 200 {object} models.LikesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/likes [get]
func GetLikesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 「いいね」の数とユーザー一覧を取得する
		rows, err := db.Query("SELECT user_id FROM likes WHERE post_id = $1", postID)
		if err != nil {
			respondError(w, "Failed to fetch likes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 「いいね」を押してくれたユーザーIDを保管する
		var userIDs []int
		for rows.Next() {
			var userID int
			if err := rows.Scan(&userID); err != nil {
				respondError(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}
			userIDs = append(userIDs, userID)
		}

		// JSONレスポンスを返す
		resp := models.LikesResponse{
			PostID:    postID,
			LikeCount: len(userIDs),
			UserIDs:   userIDs,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// UnlikePostHandler godoc
// @Summary 投稿の「いいね」を削除する
// @Description 指定したIDの投稿についている「いいね」を削除する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Produce json
// @Param id path int true "投稿ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/like [delete]
func UnlikePostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 認証情報からユーザーIDを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized", http.StatusUnauthorized)
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 「いいね」を削除する
		_, err = db.Exec("DELETE FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID)
		if err != nil {
			respondError(w, "Failed to remove like", http.StatusInternalServerError)
			return
		}

		if _, err := fmt.Fprintln(w, "like removed successfully!"); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
