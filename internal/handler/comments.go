package handler

import (
	"net/http"

	"github.com/gorilla/mux"
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
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
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
		if appErr := decodeJSON(r, &comment); appErr != nil {
			respondAppError(w, appErr)
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
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{id} [delete]
func DeleteCommentHandler(commentService *service.CommentService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからuser_idを取得
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// URIからcommentのIDを取得
		vars := mux.Vars(r)
		commentID, appErr := parseID(vars["id"])
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// コメントの所有者か確認する
		postID, err := commentService.EnsureCommentOwner(ctx, userID, commentID)
		if err != nil {
			respondAppError(w, err)
			return
		}

		// コメントを削除する
		if err := commentService.DeleteComment(ctx, commentID); err != nil {
			respondAppError(w, err)
			return
		}

		// 削除成功をJSONで返す
		respondJSON(w, http.StatusOK, map[string]string{"message": "Comment deleted successfully!"})

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
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/comments/{id} [put]
func UpdateCommentHandler(commentService *service.CommentService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// JWTからuser_idを取得
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// コメントIDを取得
		vars := mux.Vars(r)
		commentID, appErr := parseID(vars["id"])
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// コメントの所有者か確認
		postID, err := commentService.EnsureCommentOwner(ctx, userID, commentID)
		if err != nil {
			respondAppError(w, err)
			return
		}

		// リクエストボディから新しいコメント内容を取得する
		var req struct {
			Content string `json:"content"`
		}
		if appErr := decodeJSON(r, &req); appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// コメントのバリデーションを実施する
		if err := validateCommentUpdateInput(req.Content, commentID); err != nil {
			respondAppError(w, err)
			return
		}

		// コメントの更新を実施する
		if err := commentService.UpdateComment(ctx, commentID, req.Content); err != nil {
			respondAppError(w, err)
			return
		}

		// 更新成功をJSONで返す
		respondJSON(w, http.StatusOK, map[string]string{"message": "Comment update successfully!"})

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "comment_updated", UserID: userID, PostID: postID})
	}
}
