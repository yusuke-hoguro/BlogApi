package handler

import (
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"github.com/yusuke-hoguro/BlogApi/internal/service"
	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
)

// SignupHandler godoc
// @Summary 新規ユーザー登録を実施する
// @Description 送られてきたユーザー情報を使ってユーザー登録を実施する
// @Description
// @Description **エラー条件:**
// @Description - 無効なユーザー情報、ユーザー名が空、パスワードが8文字未満 → 400 Bad Request
// @Description - 許可されていないメソッド → 405 MethodNotAllowed
// @Description - データ更新/取得失敗、パスワードのハッシュ化失敗、レスポンス書き込み失敗 → 500 ServerError
// @Tags users
// @Accept json
// @Produce json
// @Param post body models.User true "ユーザー情報"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/signup [post]
func SignupHandler(userService *service.UserService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			respondAppError(w, apperror.NewAppError(apperror.TypeMethodNotAllowed, "Method Not Allowed : Method="+r.Method, nil))
			return
		}

		// GOの構造体にデコード
		var userData models.User
		if appErr := decodeJSON(r, &userData); appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// ユーザー登録のバリデーションを行う
		if err := validateSignupInput(userData); err != nil {
			respondAppError(w, err)
			return
		}

		// ユーザー登録を実施する
		if err := userService.Signup(ctx, &userData); err != nil {
			respondAppError(w, err)
			return
		}

		respondJSON(w, http.StatusCreated, userData)

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "user_signed_up", UserID: userData.ID})
	}
}

// LoginHandler godoc
// @Summary ログインする
// @Description 送られてきたユーザー情報でログインする
// @Description
// @Description **エラー条件:**
// @Description - 無効なユーザー情報、ユーザー名が空、パスワードが空 → 400 Bad Request
// @Description - ユーザー名かパスワードが不正 → 401 Unauthorized
// @Description - 許可されていないメソッド → 405 MethodNotAllowed
// @Description - データ更新/取得失敗、JWT生成失敗、レスポンス書き込み失敗 → 500 ServerError
// @Tags users
// @Accept json
// @Produce json
// @Param post body models.User true "ユーザー情報"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/login [post]
func LoginHandler(userService *service.UserService, auditPool *workerpool.AuditWorkerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストのコンテキストを取得する
		ctx := r.Context()

		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			respondAppError(w, apperror.NewAppError(apperror.TypeMethodNotAllowed, "Method Not Allowed : Method="+r.Method, nil))
			return
		}

		// GOの構造体にデコード
		var userData models.User
		if appErr := decodeJSON(r, &userData); appErr != nil {
			respondAppError(w, appErr)
			return
		}

		// ログインのバリデーションを行う
		if err := validateLoginInput(userData); err != nil {
			respondAppError(w, err)
			return
		}

		// ログインを実施する
		token, userID, err := userService.Login(ctx, userData)
		if err != nil {
			respondAppError(w, err)
			return
		}

		respondJSON(w, http.StatusOK, models.TokenResponse{Token: token})

		// 監視ワーカープールにイベントを追加
		enqueueAuditEvent(ctx, auditPool, workerpool.AuditEvent{Action: "user_logged_in", UserID: userID})
	}
}

// JWTトークンを発行する
func GenerateJWT(userID int) (string, error) {
	return service.GenerateJWT(userID)
}
