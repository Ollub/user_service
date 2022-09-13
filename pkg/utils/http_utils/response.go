package http_utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResp struct {
	Message string `json:"message"`
}

func HttpError(w http.ResponseWriter, details string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResp{details})
}

func JsonResp(w http.ResponseWriter, v interface{}, code int) {
	resp, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "Marshaling error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, "Cant write response", http.StatusInternalServerError)
	}
}
