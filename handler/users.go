package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yusuke-hoguro/BlogApi/middleware"
	"golang.org/x/crypto/bcrypt"
)

// User登録用の構造体
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ユーザー登録用のハンドラー関数
func SignupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// GOの構造体にデコード
		var userData User
		err := json.NewDecoder(r.Body).Decode(&userData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// バリデーション
		if userData.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		if len(userData.Password) < 8 {
			http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
			return
		}

		// パスワードをハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// INSERT実行
		err = db.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", userData.Username, string(hashedPassword)).Scan(&userData.ID)
		if err != nil {
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(userData)
	}
}

// ログイン機能用用のハンドラー関数
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Postであるかをチェックする
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// GOの構造体にデコード
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// バリデーション
		if creds.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		if creds.Password == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		// ユーザーをDBから検索
		var id int
		var hashedPassword string
		err = db.QueryRow("SELECT id, password FROM users WHERE username = $1", creds.Username).Scan(&id, &hashedPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		// パスワード照合
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// JWT生成
		token, err := generateJWT(id)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// レスポンス
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

// JWTトークンを発行する
func generateJWT(userID int) (string, error) {
	// payloadの生成
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), //トークンの有効時間(24時間後)
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
