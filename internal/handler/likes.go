package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/internal/service"
	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
)

// LikePostHandler godoc
// @Summary 投稿に「いいね」をつける
// @Description 指定したIDの投稿に「いいね」をつける
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags likes
// @Produce json
// @Param id path int true "投稿ID"
// @Success 201 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/like [post]
func LikePostHandler(likeService *service.LikeService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// 認証情報からユーザーIDを取得
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

		// 「いいね」を登録する
		if err := likeService.LikePost(ctx, userID, postID); err != nil {
			respondAppError(w, err)
			return
		}

		respondJSON(w, http.StatusCreated, map[string]string{"message": "Post liked successfully"})

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_liked", UserID: userID, PostID: postID})
	}
}

// GetLikesHandler godoc
// @Summary 投稿の「いいね」を取得する
// @Description 指定したIDの投稿についている「いいね」を取得する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags likes
// @Produce json
// @Param id path int true "投稿ID"
// @Success 200 {object} models.LikesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/likes [get]
func GetLikesHandler(likeService *service.LikeService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
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

		// 「いいね」の数とユーザー一覧を取得する
		likes, err := likeService.GetLikes(ctx, postID)
		if err != nil {
			respondAppError(w, err)
			return
		}

		// JSONレスポンスを返す
		respondJSON(w, http.StatusOK, likes)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "likes_fetched", PostID: postID})
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
// @Tags likes
// @Produce json
// @Param id path int true "投稿ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id}/like [delete]
func UnlikePostHandler(likeService *service.LikeService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// 認証情報からユーザーIDを取得
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

		// 「いいね」を削除する
		if err := likeService.UnlikePost(ctx, userID, postID); err != nil {
			respondAppError(w, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_unliked", UserID: userID, PostID: postID})
	}
}
