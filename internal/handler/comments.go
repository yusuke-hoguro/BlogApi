package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
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
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid post ID : PostID="+postIDStr, err))
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
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch comments : PostID="+postIDStr, err))
			return
		}
		defer rows.Close()

		// 取得したコメントをスライスに格納
		var comments []models.Comment
		for rows.Next() {
			var c models.Comment
			if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt); err != nil {
				respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Error reading comment : PostID="+postIDStr, err))
				return
			}
			comments = append(comments, c)
		}

		// rows.Next()のループが終了した後にエラーが発生していないか確認する(DBからのデータ取得中にエラーが発生していないか)
		if err := rows.Err(); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to fetch comments : PostID="+postIDStr, err))
			return
		}

		// コメントがない場合
		if len(comments) == 0 {
			fmt.Println("No comments : PostID=" + postIDStr)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comments); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
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
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid comment ID : CommentID="+IDStr, err))
			return
		}

		// 指定したIDのコメントを取得する
		var comment models.Comment
		err = db.QueryRowContext(ctx, "SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = $1", ID).Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Comment Not Found : CommentID="+IDStr, err))
			} else {
				respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Database error : CommentID="+IDStr, err))
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
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
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized", nil))
			return
		}

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postIDStr := vars["id"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid post ID : PostID="+postIDStr, err))
			return
		}

		// リクエストボディからコメントを読み取る
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body : PostID="+postIDStr, err))
			return
		}

		// コメントが空の場合はエラーとする
		if strings.TrimSpace(comment.Content) == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content is required : PostID="+postIDStr, nil))
			return
		}

		// コメントが500文字以上の場合はエラーとする
		if len(comment.Content) > MaxCommentLength {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content must be 500 characters or less : PostID="+postIDStr, nil))
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
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert comment : "+err.Error(), err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(comment); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
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
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized : User ID not found in context", nil))
			return
		}

		// URIからcommentのIDを取得
		vars := mux.Vars(r)
		commentIDStr := vars["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid comment ID : CommentID="+commentIDStr, err))
			return
		}

		// コメントの所有者か確認する
		var commentOwnerID int
		err = db.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&commentOwnerID)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Comment not found : CommentID="+commentIDStr, err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Database error : CommentID="+commentIDStr, err))
			return
		}

		// 所有者ではない場合は削除不可
		if commentOwnerID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, "Forbidden : CommentID="+commentIDStr, nil))
			return
		}

		// コメントを削除する
		_, err = db.ExecContext(ctx, "DELETE FROM comments WHERE id = $1", commentID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to delete comment : CommentID="+commentIDStr, err))
			return
		}

		// リクエスト正常終了
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, "Comment deleted successfully!"); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
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
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized : User ID not found in context", nil))
			return
		}

		// コメントIDを取得
		vars := mux.Vars(r)
		commentIDStr := vars["id"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid comment ID : CommentID="+commentIDStr, err))
			return
		}

		// リクエストボディから新しいコメント内容を取得する
		var req struct {
			Content string `json:"content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body : CommentID="+commentIDStr, err))
			return
		}

		// コメントが空の場合はエラーとする
		if strings.TrimSpace(req.Content) == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content is required : CommentID="+commentIDStr, nil))
			return
		}

		// コメントが500文字以上の場合はエラーとする
		if len(req.Content) > MaxCommentLength {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Content must be 500 characters or less : CommentID="+commentIDStr, nil))
			return
		}

		// コメントの所有者か確認
		var existringUserID int
		err = db.QueryRowContext(ctx, "SELECT user_id FROM comments WHERE id = $1", commentID).Scan(&existringUserID)
		if err == sql.ErrNoRows {
			respondAppError(w, apperror.NewAppError(apperror.TypeNotFound, "Comment not found : CommentID="+commentIDStr, err))
			return
		} else if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Database error : CommentID="+commentIDStr, err))
			return
		}
		if existringUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, "Forbidden : CommentID="+commentIDStr, nil))
			return
		}

		// コメントの更新を実施
		_, err = db.ExecContext(ctx, "UPDATE comments SET content = $1 WHERE id = $2", req.Content, commentID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to update comment : CommentID="+commentIDStr, err))
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintln(w, "Comment update successfully!"); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
			return
		}
	}
}
