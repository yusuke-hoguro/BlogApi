package models

// User はブログサービス利用者を表します。
// @Description User登録用の構造体
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse はJWTトークンを返すレスポンスを表します。
// @Description JWTトークンレスポンス用構造体
type TokenResponse struct {
	Token string `json:"token"`
}
