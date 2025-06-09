package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

// JWT認証用の秘密鍵（最終的には環境変数から）
var JwtKey = []byte("your_secret_key")

// JWTの検証を実施するミドルウェア
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// リクエストヘッダーの確認
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// JWTの解析
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 引数で指定されたハンドラー関数を実行
		next(w, r)
	}
}
