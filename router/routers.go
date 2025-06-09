package router

import (
	"database/sql"
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/handler"
	"github.com/yusuke-hoguro/BlogApi/middleware"
)

// ハンドラー関数の設定を行う
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/posts", middleware.AuthMiddleware(handler.PostHandler(db)))
	mux.HandleFunc("/posts/", middleware.AuthMiddleware(handler.PostHandler(db)))
	mux.HandleFunc("/signup", handler.SignupHandler(db))
	mux.HandleFunc("/login", handler.LoginHandler(db))
	mux.HandleFunc("/myposts", middleware.AuthMiddleware(handler.GetMyPostsHandler(db)))

}
