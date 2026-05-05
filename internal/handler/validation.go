package handler

import (
	"strings"
	"unicode/utf8"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
	"github.com/yusuke-hoguro/BlogApi/internal/models"
)

// 定数の定義
const (
	MaxTitleLength   = 100  // 投稿のタイトルの最大長
	MaxContentLength = 1000 // 投稿の内容の最大長
)

// 投稿の入力を検証する関数
func validatePostInput(post models.Post) *apperror.AppError {
	// タイトルが空の場合はエラーとする
	if strings.TrimSpace(post.Title) == "" {
		return apperror.NewAppError(apperror.TypeBadRequest, "Title must not be empty", nil)
	}

	// タイトルが100文字より大きい場合はエラーとする
	if utf8.RuneCountInString(post.Title) > MaxTitleLength {
		return apperror.NewAppError(apperror.TypeBadRequest, "Title must be 100 characters or less", nil)
	}

	// 投稿の内容が空の場合はエラーとする
	if strings.TrimSpace(post.Content) == "" {
		return apperror.NewAppError(apperror.TypeBadRequest, "Content is required", nil)
	}

	// 投稿内容が1000文字より大きい場合はエラーとする
	if utf8.RuneCountInString(post.Content) > MaxContentLength {
		return apperror.NewAppError(apperror.TypeBadRequest, "Content must be 1000 characters or less", nil)
	}
	return nil
}
