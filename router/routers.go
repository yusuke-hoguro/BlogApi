package router

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// ハンドラー関数の設定を行う
func RegisterRoutes(r *mux.Router, db *sql.DB) {
	// 投稿関係の処理
	r.HandleFunc("/api/posts", handler.GetAllPostsHandler(db)).Methods("GET")                                   // 全投稿取得用
	r.HandleFunc("/api/posts/{id}", handler.GetPostsByIDHandler(db)).Methods("GET")                             // 個別投稿取得用
	r.HandleFunc("/api/posts", middleware.AuthMiddleware(handler.CreatePostHandler(db))).Methods("POST")        // 個別投稿作成用
	r.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handler.UpdatePostHandler(db))).Methods("PUT")    // 個別投稿更新用
	r.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handler.DeletePostHandler(db))).Methods("DELETE") // 個別投稿削除用
	r.HandleFunc("/api/myposts", middleware.AuthMiddleware(handler.GetMyPostsHandler(db))).Methods("GET")       // 自身の投稿のみ取得
	// ユーザー認証系
	r.HandleFunc("/api/signup", handler.SignupHandler(db)) // ユーザー登録用
	r.HandleFunc("/api/login", handler.LoginHandler(db))   // ログイン用
	// コメント関係
	r.HandleFunc("/api/posts/{id}/comments", handler.GetCommentsByPostIDHandler(db)).Methods("GET")                     // 投稿のコメント取得
	r.HandleFunc("/api/posts/{id}/comments", middleware.AuthMiddleware(handler.PostCommentHandler(db))).Methods("POST") // 投稿のコメント投稿
	r.HandleFunc("/api/comments/{id}", handler.GetCommentsByIDHandler(db)).Methods("GET")                               // コメントIDで詳細取得
	r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(handler.DeleteCommentHandler(db))).Methods("DELETE")   // コメントIDで削除
	r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(handler.UpdateCommentHandler(db))).Methods("PUT")      // コメントを更新する
	// 「いいね」関係
	r.HandleFunc("/api/posts/{id}/like", middleware.AuthMiddleware(handler.LikePostHandler(db))).Methods("POST")     // 投稿にいいねをつける
	r.HandleFunc("/api/posts/{id}/likes", handler.GetLikesHandler(db)).Methods("GET")                                // 投稿のいいねを取得する
	r.HandleFunc("/api/posts/{id}/like", middleware.AuthMiddleware(handler.UnlikePostHandler(db))).Methods("DELETE") // 投稿のいいねを削除する
}
