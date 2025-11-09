package handler

import (
	"encoding/json"
	"net/http"
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
