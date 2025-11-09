package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusuke-hoguro/BlogApi/middleware"
	"github.com/yusuke-hoguro/BlogApi/models"
	"golang.org/x/crypto/bcrypt"
)

// ユーザー登録用のハンドラー関数
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

// ログイン機能用用のハンドラー関数
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			respondError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// GOの構造体にデコード
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			respondError(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// ユーザー名が空の場合はエラーとする
		if creds.Username == "" {
			respondError(w, "Username is required", http.StatusBadRequest)
			return
		}

		// パスワードが空の場合、エラーとする
		if creds.Password == "" {
			respondError(w, "Password is required", http.StatusBadRequest)
			return
		}

		// ユーザーをDBから検索
		var id int
		var hashedPassword string
		err = db.QueryRow("SELECT id, password FROM users WHERE username = $1", creds.Username).Scan(&id, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				respondError(w, "Invalid username or password", http.StatusUnauthorized)
			} else {
				respondError(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// パスワード照合
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
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
		if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
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
