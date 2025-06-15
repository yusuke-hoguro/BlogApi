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
	mux.HandleFunc("/posts/", routeByMethod(db))                                         //個別投稿用
	mux.HandleFunc("/signup", handler.SignupHandler(db))                                 //ユーザー登録用
	mux.HandleFunc("/login", handler.LoginHandler(db))                                   //ログイン用
	mux.HandleFunc("/myposts", middleware.AuthMiddleware(handler.GetMyPostsHandler(db))) //自身の投稿のみ取得
}

// HTTPメソッドによってJWT認証の有無判別
func routeByMethod(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetPostsyIDHandler(db)(w, r) // 認証なしのGET
			return
		}
		// 認証ありのPOST/PUT/DELETE
		middleware.AuthMiddleware(PostHandler(db))(w, r)
	}
}

// HTTPメソッドごとのルーティング
func PostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			// 記事作成用
			handler.CreatePostHandler(db)(w, r)
		case http.MethodPut:
			// 記事更新用
			handler.UpdatePostHandler(db)(w, r)
		case http.MethodDelete:
			// 記事削除
			handler.DeletePostHandler(db)(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
