package util

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, APIResponse{Code: 200, Message: "success", Data: data})
}

func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, APIResponse{Code: status, Message: msg})
}
