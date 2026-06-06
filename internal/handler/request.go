package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/middleware"
)

// ID文字列を解析して数値に変換
func parseID(idStr string) (int, *apperror.AppError) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, apperror.NewAppError(apperror.TypeBadRequest, "Invalid id: "+idStr, err)
	}
	return id, nil
}

// URLから投稿IDを抽出する関数
func postIDFromRequest(r *http.Request) (int, *apperror.AppError) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	return parseID(idStr)
}

// コンテキストからユーザーIDを取得する関数
func userIDFromContext(ctx context.Context) (int, *apperror.AppError) {
	userID, ok := ctx.Value(middleware.UserIDKey).(int)
	if !ok {
		return 0, apperror.NewAppError(apperror.TypeUnauthorized, "Unauthorized userID not found in context", nil)
	}
	return userID, nil
}

// JSONのリクエストボディを構造体にデコードする関数
func decodeJSON(r *http.Request, dst any) *apperror.AppError {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return apperror.NewAppError(apperror.TypeBadRequest, "Invalid request body", err)
	}
	return nil
}
