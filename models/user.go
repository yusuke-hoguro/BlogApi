package models

// User登録用の構造体
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// JWTトークンレスポンス用構造体
type TokenResponse struct {
	Token string `json:"token"`
}
