package models

// Like は投稿へのいいねを表します。
// @Description いいね用の構造体
type Like struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

// LikesResponse はいいね取得時のレスポンスを表します。
// @Description いいね取得時のレスポンス構造体
type LikesResponse struct {
	PostID    int   `json:"post_id"`
	LikeCount int   `json:"like_count"`
	UserIDs   []int `json:"user_ids"`
}
