package handler

import (
	"net/http"

	"backend/internal/http/middleware"
	"backend/internal/http/response"
)

type MeHandler struct{}

func NewMeHandler() MeHandler {
	return MeHandler{}
}

func (h MeHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	response.OK(w, "me", map[string]interface{}{
		"id":    claims.UserID,
		"email": claims.Email,
		"exp":   claims.Exp,
	})
}
