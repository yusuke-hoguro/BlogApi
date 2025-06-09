package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type contextKey string

// 衝突を防ぐために独自の型をキーに使用
const UserIDKey contextKey = "userID"

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

		// JWTの中身（Claims）を取り出してmap形式に変換
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// user id を保管する
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		userID := int(userIDFloat)

		// ユーザーIDをリクエストのContextに埋め込んで次のハンドラー関数に渡す
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		// 引数で指定されたハンドラー関数を実行
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
