package handler

import (
	"encoding/json"
	"net/http"
)

// JSONレスポンスを返す共通関数
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// エラーレスポンスを返す関数
func respondError(w http.ResponseWriter, staus int, message string) {
	respondJSON(w, staus, map[string]string{"message": message})
}
