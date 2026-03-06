package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
)

const (
	MaxCommentLength = 500
)

// GetCommentsByPostIDHandler godoc
// @Summary 投稿のコメントを取得する
// @Description 指定した投稿のコメントをすべて取得する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Produce json
// @Param id path int true "投稿ID"
// @Success 200 {array} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/comments [get]
func GetCommentsByPostIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID : PostID="+postIDStr, http.StatusBadRequest)
			return
		}

		// 指定した投稿のコメントを取得する
		rows, err := db.QueryContext(ctx, `
			SELECT id, post_id, user_id, content, created_at
			FROM comments
			WHERE post_id = $1
			ORDER BY created_at ASC
		`, postID)
		if err != nil {
			respondError(w, "Failed to fetch comments : PostID="+postIDStr, http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// 取得したコメントをスライスに格納
		var comments []models.Comment
		for rows.Next() {
			var c models.Comment
			if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
				respondError(w, "Error reading comment : PostID="+postIDStr, http.StatusInternalServerError)
				return
			}
			comments = append(comments, c)
		}

		// rows.Next()のループが終了した後にエラーが発生していないか確認する(DBからのデータ取得中にエラーが発生していないか)
		if err := rows.Err(); err != nil {
			respondError(w, "Failed to fetch comments : PostID="+postIDStr, http.StatusInternalServerError)
			return
		}

		// コメントがない場合
		if len(comments) == 0 {
			fmt.Println("No comments : PostID=" + postIDStr)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GetCommentsByIDHandler godoc
// @Summary 指定したコメントを取得する
// @Description コメントIDを指定してコメントを取得する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - コメントが存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Produce json
// @Param id path int true "コメントID"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{id} [get]
func GetCommentsByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// URIからコメントのIDを取得
		vars := mux.Vars(r)
		IDStr := vars["id"]
		ID, err := strconv.Atoi(IDStr)
		if err != nil {
			respondError(w, "Invalid comment ID : CommentID="+IDStr, http.StatusBadRequest)
			return
		}

		// 指定したIDのコメントを取得する
		var comment models.Comment
		err = db.QueryRowContext(ctx, "SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1", ID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				respondError(w, "Comment Not Found : CommentID="+IDStr, http.StatusNotFound)
			} else {
				respondError(w, "Database error : CommentID="+IDStr, http.StatusInternalServerError)
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

// PostCommentHandler godoc
// @Summary 指定した投稿にコメントを追加する
// @Description 指定された投稿に送られてきたコメントを追加する
// @Description
// @Description **エラー条件:**
// @Description - 無効な投稿ID、無効なコメント内容、コメントが空、コメントが500文字以上 → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "投稿ID"
// @Param post body models.Comment true "コメント内容"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/comments [post]
func PostCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからuser_idを取得
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondError(w, "Invalid post ID : PostID="+postIDStr, http.StatusBadRequest)
			return
		}

		// リクエストボディからコメントを読み取る
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			respondError(w, "Invalid request body : PostID="+postIDStr, http.StatusBadRequest)
			return
		}

		// コメントが空の場合はエラーとする
		if strings.TrimSpace(comment.Content) == "" {
			respondError(w, "Content is required : PostID="+postIDStr, http.StatusBadRequest)
		}

		// コメントが500文字以上の場合はエラーとする
		if len(comment.Content) > MaxCommentLength {
			respondError(w, "Content must be 500 characters or less : PostID="+postIDStr, http.StatusBadRequest)
			return
		}

		// コメントを挿入する
		query := `INSERT INTO comments (post_id, user_id, content) 
				VALUES ($1, $2, $3)
				RETURNING id, post_id, user_id, content, created_at`

		err = db.QueryRowContext(ctx, query, postID, userID, comment.Content).Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			respondError(w, "Failed to insert comment : "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// DeleteCommentHandler godoc
// @Summary 指定したコメントを削除する
// @Description 送られてきたIDのコメントを削除する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 送信者がコメントの所有者でない → 403 Forbidden を返す
// @Description - コメントが存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "コメントID"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{id} [delete]
func DeleteCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからuser_idを取得
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized : User ID not found in context", http.StatusUnauthorized)
			return
		}

		// URIからcommentのIDを取得
		vars := mux.Vars(r)
		commentIDStr := vars["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			respondError(w, "Invalid comment ID : CommentID="+commentIDStr, http.StatusBadRequest)
			return
		}

		// コメントの所有者か確認する
		var commentOwnerID int
		err = db.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&commentOwnerID)
		if err == sql.ErrNoRows {
			respondError(w, "Comment not found", http.StatusNotFound)
			return
		} else if err != nil {
			respondError(w, "Database error : CommentID="+commentIDStr, http.StatusInternalServerError)
			return
		}

		// 所有者ではない場合は削除不可
		if commentOwnerID != userID {
			respondError(w, "Forbidden : CommentID="+commentIDStr, http.StatusForbidden)
			return
		}

		// コメントを削除する
		_, err = db.ExecContext(ctx, "DELETE FROM comments WHERE id = $1", commentID)
		if err != nil {
			respondError(w, "Failed to delete comment : CommentID="+commentIDStr, http.StatusInternalServerError)
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

// UpdateCommentHandler godoc
// @Summary コメントの内容を更新する
// @Description 送られてきたIDのコメントを更新する。
// @Description
// @Description **エラー条件:**
// @Description - 無効なID、空コメント、500文字以上のコメント、コメント未取得 → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 送信者がコメントの所有者でない → 403 Forbidden
// @Description - コメントが存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags comments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "コメントID"
// @Param post body models.Comment true "コメント内容"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{id} [put]
func UpdateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからuser_idを取得
		userID, ok := ctx.Value(middleware.UserIDKey).(int)
		if !ok {
			respondError(w, "Unauthorized : User ID not found in context", http.StatusUnauthorized)
			return
		}

		// コメントIDを取得
		vars := mux.Vars(r)
		commentIDStr := vars["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			respondError(w, "Invalid comment ID : CommentID="+commentIDStr, http.StatusBadRequest)
			return
		}

		// リクエストボディから新しいコメント内容を取得する
		var req struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondError(w, "Invalid request body : CommentID="+commentIDStr, http.StatusBadRequest)
			return
		}

		// コメントが空の場合はエラーとする
		if strings.TrimSpace(req.Content) == "" {
			respondError(w, "Content is required : CommentID="+commentIDStr, http.StatusBadRequest)
		}

		// コメントが500文字以上の場合はエラーとする
		if len(req.Content) > MaxCommentLength {
			respondError(w, "Content must be 500 characters or less : CommentID="+commentIDStr, http.StatusBadRequest)
			return
		}

		// コメントの所有者か確認
		var existringUserID int
		err = db.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&existringUserID)
		if err == sql.ErrNoRows {
			respondError(w, "Comment not found : CommentID="+commentIDStr, http.StatusNotFound)
			return
		} else if err != nil {
			respondError(w, "Database error : CommentID="+commentIDStr, http.StatusInternalServerError)
			return
		}
		if existringUserID != userID {
			respondError(w, "Forbidden : CommentID="+commentIDStr, http.StatusForbidden)
			return
		}

		// コメントの更新を実施
		_, err = db.ExecContext(ctx, "UPDATE comments SET content = $1 WHERE id = $2", req.Content, commentID)
		if err != nil {
			respondError(w, "Failed to update comment : CommentID="+commentIDStr, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, "Comment update successfully!"); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
