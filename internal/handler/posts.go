package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/repository"
	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
)

// GetPostsByIDHandler godoc
// @Summary 投稿をIDで取得する
// @Description 指定したIDの投稿を返す
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Produce json
// @Param id path int true "PostID"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id} [get]
func GetPostsByIDHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()
		// URLから投稿IDを抽出する
		id, appErr := postIDFromRequest(r)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// DBから指定したIDの投稿を取得する
		repo := repository.NewPostRepository(db)
		post, err := repo.FindByID(ctx, id)
		if err != nil {
			respondAppError(w, err)
			return
		}
		// 取得した投稿をJSONで返す
		respondJSON(w, http.StatusOK, post)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_fetched", UserID: post.UserID, PostID: post.ID})
	}
}

// CreatePostHandler godoc
// @Summary 新しい投稿を作成する
// @Description 送られてきた構造体のデータから新規投稿を作成する
// @Description
// @Description **エラー条件:**
// @Description - 無効な投稿内容、タイトルか投稿内容が空、タイトルが100文字以上、投稿内容が1000文字以上 → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param post body models.Post true "投稿内容"
// @Success 201 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts [post]
func CreatePostHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()
		// JWTからユーザーIDを取得する
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// Post型の構造体にデコードして格納
		var post models.Post
		if appErr := decodeJSON(r, &post); appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// 投稿のバリデーションを行う
		if err := validatePostInput(post); err != nil {
			respondAppError(w, err)
			return
		}
		// 記事にユーザーIDを設定する
		post.UserID = userID
		// 投稿を作成する
		repo := repository.NewPostRepository(db)
		if err := repo.Create(ctx, &post); err != nil {
			respondAppError(w, err)
			return
		}
		// 作成した投稿をJSONで返す
		respondJSON(w, http.StatusCreated, post)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_created", UserID: post.UserID, PostID: post.ID})
	}
}

// UpdatePostHandler godoc
// @Summary 投稿の内容を更新する
// @Description 送られてきた構造体のデータから投稿を更新する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID、無効な投稿内容、タイトルか投稿内容が空、タイトルが100文字以上、投稿内容が1000文字以上 → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 送信者が記事の投稿者でない → 403 Forbidden
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "投稿ID"
// @Param post body models.Post true "投稿内容"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id} [put]
func UpdatePostHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()
		// JWTからユーザーIDを取得する
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// URLから投稿IDを抽出する
		id, appErr := postIDFromRequest(r)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// DBから投稿者のユーザーIDを取得する
		repo := repository.NewPostRepository(db)
		postUserID, err := repo.FindUserIDByPostID(ctx, id)
		if err != nil {
			respondAppError(w, err)
			return
		}
		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, fmt.Sprintf("Forbidden : PostID=%d", id), nil))
			return
		}
		// Post型の構造体にデコードして格納
		var post models.Post
		if appErr := decodeJSON(r, &post); appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// 投稿のバリデーションを行う
		if err := validatePostInput(post); err != nil {
			respondAppError(w, err)
			return
		}
		// 指定したIDの投稿を更新する
		if err := repo.Update(ctx, id, &post); err != nil {
			respondAppError(w, err)
			return
		}
		// 更新した投稿をJSONで返す
		respondJSON(w, http.StatusOK, post)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_updated", UserID: userID, PostID: id})
	}
}

// DeletePostHandler godoc
// @Summary 投稿を削除する
// @Description 送られてきたIDの投稿を削除する
// @Description
// @Description **エラー条件:**
// @Description - 無効なID → 400 Bad Request
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 送信者が記事の投稿者でない → 403 Forbidden
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "投稿ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts/{id} [put]
func DeletePostHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()
		// JWTからリクエストをなげたユーザーIDを取得
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// URLからIDを取得する
		id, appErr := postIDFromRequest(r)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// 削除対象の投稿を作成したユーザーのIDを取得する
		repo := repository.NewPostRepository(db)
		postUserID, err := repo.FindUserIDByPostID(ctx, id)
		if err != nil {
			respondAppError(w, err)
			return
		}
		// リクエストを投げたユーザーが記事の投稿者でない場合はエラー
		if postUserID != userID {
			respondAppError(w, apperror.NewAppError(apperror.TypeForbidden, fmt.Sprintf("Forbidden : PostID=%d", id), nil))
			return
		}
		// DELETE実行
		if err := repo.Delete(ctx, id); err != nil {
			respondAppError(w, err)
			return
		}
		// 削除成功のため204 No Contentを返す
		respondJSON(w, http.StatusNoContent, nil)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "post_deleted", UserID: userID, PostID: id})
	}
}

// GetMyPostsHandler godoc
// @Summary ユーザー自身の投稿を取得する
// @Description リクエストを投げたユーザーが作成した投稿を取得する
// @Description
// @Description **エラー条件:**
// @Description - リクエスト認証エラー → 401 Unauthorized
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {array} models.Post
// @Failure 401 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/myposts [get]
func GetMyPostsHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()
		// JWTからリクエストをなげたユーザーIDを取得
		userID, appErr := userIDFromContext(ctx)
		if appErr != nil {
			respondAppError(w, appErr)
			return
		}
		// DBから指定したユーザーIDの投稿を取得する
		repo := repository.NewPostRepository(db)
		posts, err := repo.ListByUserID(ctx, userID)
		if err != nil {
			respondAppError(w, err)
			return
		}
		// 取得した投稿をJSONで返す
		respondJSON(w, http.StatusOK, posts)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "my_posts_fetched", UserID: userID})
	}
}

// GetAllPostsHandler godoc
// @Summary すべての投稿を取得する
// @Description DBから全投稿を取得して返却する
// @Description 送られてきたIDの投稿を削除する
// @Description
// @Description **エラー条件:**
// @Description - 投稿が存在しない → 404 Not Found
// @Description - データ更新/取得失敗 or レスポンス書き込み失敗 → 500 ServerError
// @Tags posts
// @Produce json
// @Success 200 {array} models.Post
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/posts [get]
func GetAllPostsHandler(db *sql.DB, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()
		// 全投稿を取得する
		repo := repository.NewPostRepository(db)
		posts, err := repo.ListAll(ctx)
		if err != nil {
			respondAppError(w, err)
			return
		}
		// 取得した投稿をJSONで返す
		respondJSON(w, http.StatusOK, posts)
		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "posts_fetched"})
	}
}
