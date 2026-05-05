package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
)

func postIDFromRequest(r *http.Request) (int, *apperror.AppError) {
	// URLからIDを取得する
	idStr := strings.TrimPrefix(r.URL.Path, "/api/posts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, apperror.NewAppError(apperror.TypeBadRequest, "Invalid ID : PostID="+idStr, err)
	}
	return id, nil
}
