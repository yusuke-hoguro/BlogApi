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
	"github.com/yusuke-hoguro/BlogApi/internal/service"
	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
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
func GetCommentsByPostIDHandler(commentService *service.CommentService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// URIからpostのIDを取得
		vars := mux.Vars(r)
		postID, appErr := parseID(vars["id"])
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// 指定した投稿のコメントをすべて取得する
		comments, err := commentService.GetCommentsByPostID(ctx, postID)
		if err != nil {
			respondAppError(w, err)
			return
		}

		// 指定した投稿のコメントをJSONで返す
		respondJSON(w, http.StatusOK, comments)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "comments_fetched", PostID: postID})
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
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{id} [get]
func GetCommentsByIDHandler(commentService *service.CommentService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// URIからコメントのIDを取得
		vars := mux.Vars(r)
		id, appErr := parseID(vars["id"])
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// 指定したIDのコメントを取得する
		comment, err := commentService.GetCommentByID(ctx, id)
		if err != nil {
			respondAppError(w, err)
			return
		}

		// 指定したコメントをJSONで返す
		respondJSON(w, http.StatusOK, comment)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "comment_fetched", UserID: comment.UserID, PostID: comment.PostID})
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
// @Success 201 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/comments [post]
func PostCommentHandler(commentService *service.CommentService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
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
		postID, appErr := parseID(vars["id"])
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// リクエストボディからコメントを読み取る
		var comment models.Comment
		if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, fmt.Sprintf("Invalid request body : PostID=%d", postID), err))
			return
		}

		// コメントのバリデーションを実施する
		if err := validateCommentInput(comment, postID); err != nil {
			respondAppError(w, err)
			return
		}

		// コメントを挿入する
		if err := commentService.CreateComment(ctx, postID, userID, &comment); err != nil {
			respondAppError(w, err)
			return
		}

		// 作成したコメントをJSONで返す
		respondJSON(w, http.StatusCreated, comment)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "comment_created", UserID: comment.UserID, PostID: comment.PostID})
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
func DeleteCommentHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
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
		var commentOwnerID, postID int
		err = db.QueryRowContext(ctx, "SELECT user_id, post_id FROM comments WHERE id = $1", commentID).Scan(&commentOwnerID, &postID)
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

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "comment_deleted", UserID: userID, PostID: postID})
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
func UpdateCommentHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
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
		var existringUserID, postID int
		err = db.QueryRowContext(ctx, "SELECT user_id, post_id FROM comments WHERE id = $1", commentID).Scan(&existringUserID, &postID)
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

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "comment_updated", UserID: userID, PostID: postID})
	}
}
