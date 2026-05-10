package router

import (
	"database/sql"
	"net/http"

	_ "github.com/yusuke-hoguro/BlogApi/docs" // swag init で生成されたdocsパッケージをimport

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/yusuke-hoguro/BlogApi/internal/app"
	"github.com/yusuke-hoguro/BlogApi/internal/handler"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
	"github.com/yusuke-hoguro/BlogApi/internal/workerpool"
)

// ハンドラー関数の設定を行う
func RegisterRoutes(r *mux.Router, db *sql.DB, auditPool *workerpool.AuditWorkerPool, services *app.Services) {
	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// ヘルスチェック用
	r.HandleFunc("/api/healthz", handler.HealthzHandler(auditPool)).Methods(http.MethodGet, http.MethodHead) // ヘルスチェック用
	// 投稿関係の処理
	r.HandleFunc("/api/posts", handler.GetAllPostsHandler(services.Post, auditPool)).Methods(http.MethodGet)                                   // 全投稿取得用
	r.HandleFunc("/api/posts/{id}", handler.GetPostsByIDHandler(services.Post, auditPool)).Methods(http.MethodGet)                             // 個別投稿取得用
	r.HandleFunc("/api/posts", middleware.AuthMiddleware(handler.CreatePostHandler(services.Post, auditPool))).Methods(http.MethodPost)        // 個別投稿作成用
	r.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handler.UpdatePostHandler(services.Post, auditPool))).Methods(http.MethodPut)    // 個別投稿更新用
	r.HandleFunc("/api/posts/{id}", middleware.AuthMiddleware(handler.DeletePostHandler(services.Post, auditPool))).Methods(http.MethodDelete) // 個別投稿削除用
	r.HandleFunc("/api/myposts", middleware.AuthMiddleware(handler.GetMyPostsHandler(services.Post, auditPool))).Methods(http.MethodGet)       // 自身の投稿のみ取得
	// ユーザー認証系
	r.HandleFunc("/api/signup", handler.SignupHandler(db, auditPool)).Methods(http.MethodPost) // ユーザー登録用
	r.HandleFunc("/api/login", handler.LoginHandler(db, auditPool)).Methods(http.MethodPost)   // ログイン用
	// コメント関係
	r.HandleFunc("/api/posts/{id}/comments", handler.GetCommentsByPostIDHandler(db, auditPool)).Methods(http.MethodGet)                     // 投稿のコメント取得
	r.HandleFunc("/api/posts/{id}/comments", middleware.AuthMiddleware(handler.PostCommentHandler(db, auditPool))).Methods(http.MethodPost) // 投稿のコメント投稿
	r.HandleFunc("/api/comments/{id}", handler.GetCommentsByIDHandler(db, auditPool)).Methods(http.MethodGet)                               // コメントIDで詳細取得
	r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(handler.DeleteCommentHandler(db, auditPool))).Methods(http.MethodDelete)   // コメントIDで削除
	r.HandleFunc("/api/comments/{id}", middleware.AuthMiddleware(handler.UpdateCommentHandler(db, auditPool))).Methods(http.MethodPut)      // コメントを更新する
	// 「いいね」関係
	r.HandleFunc("/api/posts/{id}/like", middleware.AuthMiddleware(handler.LikePostHandler(db, auditPool))).Methods(http.MethodPost)     // 投稿にいいねをつける
	r.HandleFunc("/api/posts/{id}/likes", handler.GetLikesHandler(db, auditPool)).Methods(http.MethodGet)                                // 投稿のいいねを取得する
	r.HandleFunc("/api/posts/{id}/like", middleware.AuthMiddleware(handler.UnlikePostHandler(db, auditPool))).Methods(http.MethodDelete) // 投稿のいいねを削除する
}
