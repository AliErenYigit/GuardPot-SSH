package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"backend/internal/domain"
	"backend/internal/http/middleware"
	"backend/internal/http/response"
	"backend/internal/http/dto"
	"backend/internal/repository"
	"backend/internal/service"
)

type SSHHandler struct {
	ssh service.SSHService
}

func NewSSHHandler(ssh service.SSHService) SSHHandler {
	return SSHHandler{ssh: ssh}
}

func (h SSHHandler) Create(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	var req dto.CreateSSHConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, http.StatusBadRequest, "bad_request", "Invalid JSON body")
		return
	}

	authType := domain.SSHAuthType(req.AuthType)
	item, err := h.ssh.Create(r.Context(), auth.UserID, req.Name, req.Host, req.Port, req.Username, authType, req.Secret)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "invalid_input", err.Error())
		return
	}

	response.Created(w, "ssh connection created", item)
}

func (h SSHHandler) List(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	items, err := h.ssh.List(r.Context(), auth.UserID)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		return
	}

	response.OK(w, "ssh connections", items)
}

func (h SSHHandler) Delete(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(w, http.StatusBadRequest, "invalid_id", "Invalid id")
		return
	}

	err = h.ssh.Delete(r.Context(), auth.UserID, id)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Fail(w, http.StatusNotFound, "not_found", "Connection not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		return
	}

	response.OK(w, "deleted", map[string]any{"id": id})
}
