package models

import "time"

// Post はブログ投稿を表します。
// @Description ブログ投稿用の構造体
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
