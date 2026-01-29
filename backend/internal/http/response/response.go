package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func OK(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, SuccessResponse{Message: message, Data: data})
}

func Created(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusCreated, SuccessResponse{Message: message, Data: data})
}

func Fail(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}
