package handler

import (
	"encoding/json"
	"net/http"

	"backend/internal/http/dto"
	"backend/internal/http/response"
	"backend/internal/service"
)

type AuthHandler struct {
	auth service.AuthService
}

func NewAuthHandler(auth service.AuthService) AuthHandler {
	return AuthHandler{auth: auth}
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "bad_request", "Invalid JSON body")
		return
	}

	user, err := h.auth.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidEmail:
			response.Fail(w, http.StatusBadRequest, "invalid_email", "Email is invalid")
		case service.ErrInvalidPassword:
			response.Fail(w, http.StatusBadRequest, "invalid_password", "Password must be at least 8 characters")
		case service.ErrEmailAlreadyInUse:
			response.Fail(w, http.StatusConflict, "email_in_use", "Email is already in use")
		default:
			response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		}
		return
	}

	response.Created(w, "user created", user)
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "bad_request", "Invalid JSON body")
		return
	}

	res, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrInvalidCredentials:
			response.Fail(w, http.StatusUnauthorized, "invalid_credentials", "Email or password is incorrect")
		default:
			response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		}
		return
	}

	response.OK(w, "login successful", res)
}
