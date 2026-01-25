package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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
		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			respondError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// GOの構造体にデコード
		var userData models.User
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// ユーザー名が空の場合はエラーとする
		if userData.Username == "" {
			respondError(w, "Username is required", http.StatusBadRequest)
			return
		}

		// パスワードが8文字未満の場合はエラーとする
		if len(userData.Password) < 8 {
			respondError(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}

		// パスワードをハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			respondError(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// INSERT実行
		err = db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", userData.Username, string(hashedPassword)).Scan(&userData.ID)
		if err != nil {
			respondError(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(userData); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
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
		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			respondError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// GOの構造体にデコード
		var userData models.User
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// ユーザー名が空の場合はエラーとする
		if userData.Username == "" {
			respondError(w, "Username is required", http.StatusBadRequest)
			return
		}

		// パスワードが空の場合、エラーとする
		if userData.Password == "" {
			respondError(w, "Password is required", http.StatusBadRequest)
			return
		}

		// ユーザーをDBから検索
		var id int
		var hashedPassword string
		err = db.QueryRow("SELECT id, password FROM users WHERE username = $1", userData.Username).Scan(&id, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				respondError(w, "Invalid username or password", http.StatusUnauthorized)
			} else {
				respondError(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// パスワード照合
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userData.Password))
		if err != nil {
			respondError(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// JWT生成
		token, err := GenerateJWT(id)
		if err != nil {
			respondError(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// レスポンス
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(models.TokenResponse{Token: token}); err != nil {
			respondError(w, err.Error(), http.StatusInternalServerError)
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
