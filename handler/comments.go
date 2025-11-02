package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/middleware"
	"github.com/yusuke-hoguro/BlogApi/models"
)

const (
	MaxCommentLength = 500
)

// 投稿のコメント取得用ハンドラー関数
func GetCommentsByPostIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID", http.StatusBadRequest)
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
			respondError(w, "Failed to fetch comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 取得したコメントをスライスに格納
		var comments []models.Comment
		for rows.Next() {
			var c models.Comment
			if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
				respondError(w, "Error reading comment", http.StatusInternalServerError)
				return
			}
			comments = append(comments, c)
		}

		// コメントがない場合
		if len(comments) == 0 {
			fmt.Println("No comments.")
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// 投稿のコメント取得用ハンドラー関数
func GetCommentsByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URIからpostのIDを取得
		vars := mux.Vars(r)
		IDStr := vars["id"]
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			respondError(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// 指定したIDのコメントを取得する
		var comment models.Comment
		err = db.QueryRow("SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1", ID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				respondError(w, "Comment Not Found", http.StatusInternalServerError)
			} else {
				respondError(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// コメント投稿用のハンドラー関数
func PostCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからuser_idを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID", http.StatusBadRequest)
			return
		}

		// リクエストボディからコメントを読み取る
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// コメントが空の場合はエラーとする
		if strings.TrimSpace(comment.Content) == "" {
			respondError(w, "Content is required", http.StatusBadRequest)
		}

		// コメントが500文字以上の場合はエラーとする
		if len(comment.Content) > MaxCommentLength {
			respondError(w, "Content must be 500 characters or less", http.StatusBadRequest)
			return
		}

		// コメントを挿入する
		query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3)`
		_, err = db.Exec(query, postID, userID, comment.Content)
		if err != nil {
			respondError(w, "Failed to insert comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]string{"message": "Comment created"}); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// コメント削除用のハンドラー関数
func DeleteCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからuser_idを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// URIからcommentのIDを取得
		vars := mux.Vars(r)
		commentIDStr := vars["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			respondError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// コメントの所有者か確認する
		var commentOwnerID int
		err = db.QueryRow("SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&commentOwnerID)
		if err == sql.ErrNoRows {
			respondError(w, "Comment not found", http.StatusNotFound)
			return
		} else if err != nil {
			respondError(w, "Database error", http.StatusInternalServerError)
			return
		}

		// 所有者ではない場合は削除不可
		if commentOwnerID != userID {
			respondError(w, "Forbidden", http.StatusForbidden)
			return
		}

		// コメントを削除する
		_, err = db.Exec("DELETE FROM comments WHERE id = $1", commentID)
		if err != nil {
			respondError(w, "Failed to delete comment", http.StatusInternalServerError)
			return
		}

		// リクエスト正常終了
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, "Comment deleted successfully!"); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// コメント更新用のハンドラー関数
func UpdateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// JWTからuser_idを取得
		userID, ok := r.Context().Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// コメントIDを取得
		vars := mux.Vars(r)
		commentIDStr := vars["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			respondError(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}

		// リクエストボディから新しいコメント内容を取得する
		var req struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// コメントが空の場合はエラーとする
		if strings.TrimSpace(req.Content) == "" {
			respondError(w, "Content is required", http.StatusBadRequest)
		}

		// コメントが500文字以上の場合はエラーとする
		if len(req.Content) > MaxCommentLength {
			respondError(w, "Content must be 500 characters or less", http.StatusBadRequest)
			return
		}

		// コメントの所有者か確認
		var existringUserID int
		err = db.QueryRow("SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&existringUserID)
		if err == sql.ErrNoRows {
			respondError(w, "Comment not found", http.StatusNotFound)
			return
		} else if err != nil {
			respondError(w, "Database error", http.StatusInternalServerError)
			return
		}
		if existringUserID != userID {
			respondError(w, "Forbidden", http.StatusForbidden)
			return
		}

		// コメントの更新を実施
		_, err = db.Exec("UPDATE comments SET content = $1 WHERE id = $2", req.Content, commentID)
		if err != nil {
			respondError(w, "Failed to update comment", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, "Comment update successfully!"); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
