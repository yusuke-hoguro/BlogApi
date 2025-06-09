package router

import (
	"database/sql"
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// ハンドラー関数の設定を行う
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/posts", middleware.AuthMiddleware(handler.GetPostsHandler(db)))
	mux.HandleFunc("/posts/", middleware.AuthMiddleware(handler.CreatePostHandler(db)))
	mux.HandleFunc("/signup", handler.SignupHandler(db))
	mux.HandleFunc("/login", handler.LoginHandler(db))

	//続きはここから。ハンドラー関数も移動必要
}
