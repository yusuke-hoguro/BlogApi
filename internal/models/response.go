package models

// ErrorResponse はエラーレスポンスを表します。
// @Description エラーレスポンス構造体
type ErrorResponse struct {
	Message string `json:"message"`
}
