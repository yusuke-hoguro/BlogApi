package models

type Like struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	PostID int `json:"post_id"`
}

// いいね取得時のレスポンス用
type LikesResponse struct {
	PostID    int   `json:"post_id"`
	LikeCount int   `json:"like_count"`
	UserIDs   []int `json:"user_ids"`
}
