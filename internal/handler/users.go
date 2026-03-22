package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
	"golang.org/x/crypto/bcrypt"
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
func SignupHandler(db *sql.DB) http.HandlerFunc {
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
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body", err))
			return
		}

		// ユーザー名が空の場合はエラーとする
		if userData.Username == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Username is required", nil))
			return
		}

		// パスワードが8文字未満の場合はエラーとする
		if len(userData.Password) < 8 {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Password must be at least 8 characters long", nil))
			return
		}

		// パスワードをハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to hash password : Username="+userData.Username, err))
			return
		}

		// INSERT実行
		err = db.QueryRowContext(ctx, "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", userData.Username, string(hashedPassword)).Scan(&userData.ID)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to insert user : Username="+userData.Username, err))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
			return
		}
	}
}

// LoginHandler godoc
// @Summary ログインする
// @Description 送られてきたユーザー情報でログインする
// @Description
// @Description **エラー条件:**
// @Description - 無効なユーザー情報、ユーザー名が空、パスワードが8文字未満 → 400 Bad Request
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
func LoginHandler(db *sql.DB) http.HandlerFunc {
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
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body", err))
			return
		}

		// ユーザー名が空の場合はエラーとする
		if userData.Username == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Username is required", nil))
			return
		}

		// パスワードが空の場合、エラーとする
		if userData.Password == "" {
			respondAppError(w, apperror.NewAppError(apperror.TypeBadRequest, "Password is required : Username="+userData.Username, nil))
			return
		}

		// ユーザーをDBから検索
		var id int
		var hashedPassword string
		err = db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE username = $1", userData.Username).Scan(&id, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Invalid username or password : Username="+userData.Username, err))
			} else {
				respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Database error : Username="+userData.Username, err))
			}
			return
		}

		// パスワード照合
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userData.Password))
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeUnauthorized, "Invalid username or password : Username="+userData.Username, err))
			return
		}

		// JWT生成
		token, err := GenerateJWT(id)
		if err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to generate token : Username="+userData.Username, err))
			return
		}

		// レスポンス
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models.TokenResponse{Token: token}); err != nil {
			respondAppError(w, apperror.NewAppError(apperror.TypeInternalServer, "Failed to write response", err))
			return
		}
	}
}

// JWTトークンを発行する
func GenerateJWT(userID int) (string, error) {
	// payloadの生成
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // トークンの有効時間(24時間後)
	}

	// JWTを生成する
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 署名付きトークン生成
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
