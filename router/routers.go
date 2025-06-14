package router

import (
	"database/sql"
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// ハンドラー関数の設定を行う
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/posts", handler.GetAllPostsHandler(db))                             //全投稿用
	mux.HandleFunc("/posts/", middleware.AuthMiddleware(handler.PostHandler(db)))        //個別投稿用
	mux.HandleFunc("/signup", handler.SignupHandler(db))                                 //ユーザー登録用
	mux.HandleFunc("/login", handler.LoginHandler(db))                                   //ログイン用
	mux.HandleFunc("/myposts", middleware.AuthMiddleware(handler.GetMyPostsHandler(db))) //自身の投稿のみ取得
}
