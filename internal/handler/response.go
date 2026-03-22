package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/yusuke-hoguro/BlogApi/internal/apperror"
)

// JSONレスポンスを返す共通関数
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// エラーレスポンスを返す関数
func respondError(w http.ResponseWriter, message string, staus int) {
	respondJSON(w, staus, map[string]string{"message": message})
}

// アプリケーションエラーを処理する関数
func respondAppError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	// エラーの中にAppError構造体が含まれているか確認
	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			log.Printf("app error: type=%s message=%s cause=%v", appErr.Type, appErr.Message, appErr.Err)
		} else {
			log.Printf("app error: type=%s message=%s", appErr.Type, appErr.Message)
		}
		respondError(w, appErr.Message, apperror.GetStatusCode(appErr.Type))
		return
	}
	log.Printf("unexpected error: %v", err)
	respondError(w, "Internal Server Error", http.StatusInternalServerError)
}
