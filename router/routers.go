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
	r.HandleFunc("/posts", handler.GetAllPostsHandler(db)).Methods("GET")                                   //全投稿取得用
	r.HandleFunc("/posts/{id}", handler.GetPostsByIDHandler(db)).Methods("GET")                             //個別投稿取得用
	r.HandleFunc("/posts", middleware.AuthMiddleware(handler.CreatePostHandler(db))).Methods("POST")        //個別投稿作成用
	r.HandleFunc("/posts/{id}", middleware.AuthMiddleware(handler.UpdatePostHandler(db))).Methods("PUT")    //個別投稿更新用
	r.HandleFunc("/posts/{id}", middleware.AuthMiddleware(handler.DeletePostHandler(db))).Methods("DELETE") //個別投稿削除用
	r.HandleFunc("/myposts", middleware.AuthMiddleware(handler.GetMyPostsHandler(db))).Methods("GET")       //自身の投稿のみ取得
	// ユーザー認証系
	r.HandleFunc("/signup", handler.SignupHandler(db)) //ユーザー登録用
	r.HandleFunc("/login", handler.LoginHandler(db))   //ログイン用
	// コメント取得
	r.HandleFunc("/posts/{id}/comments", handler.GetCommentsByPostIDHandler(db)).Methods("GET") //投稿のコメント取得
}
