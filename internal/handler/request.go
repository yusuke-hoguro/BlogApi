package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
)

// URLから投稿IDを抽出する関数
func postIDFromRequest(r *http.Request) (int, *apperror.AppError) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, apperror.NewAppError(apperror.TypeBadRequest, "Invalid ID : PostID="+idStr, err)
	}
	return id, nil
}

// コンテキストからユーザーIDを取得する関数
func userIDFromContext(ctx context.Context) (int, *apperror.AppError) {
	userID, ok := ctx.Value(middleware.UserIDKey).(int)
	if !ok {
		return 0, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized userID not found in context", nil)
	}
	return userID, nil
}
