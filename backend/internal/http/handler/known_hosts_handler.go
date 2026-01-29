package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"backend/internal/http/middleware"
	"backend/internal/http/response"
	"backend/internal/repository"
	"backend/internal/service"
)

type KnownHostsHandler struct {
	ssh   service.SSHService
	kh    service.KnownHostsManager
}

func NewKnownHostsHandler(ssh service.SSHService, kh service.KnownHostsManager) KnownHostsHandler {
	return KnownHostsHandler{ssh: ssh, kh: kh}
}

// GET /api/v1/ssh/known-hosts
func (h KnownHostsHandler) List(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	lines, err := h.kh.ListLines()
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, "server_error", "Cannot read known_hosts")
		return
	}
	response.OK(w, "known_hosts", lines)
}

// POST /api/v1/ssh/{id}/hostkey/scan
func (h KnownHostsHandler) Scan(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id <= 0 {
		response.Fail(w, http.StatusBadRequest, "invalid_id", "Invalid id")
		return
	}

	meta, _, err := h.ssh.GetDecrypted(r.Context(), auth.UserID, id)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Fail(w, http.StatusNotFound, "not_found", "Connection not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		return
	}

	res, err := service.ScanHostKey(r.Context(), meta.Host, meta.Port)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "scan_failed", err.Error())
		return
	}

	response.OK(w, "scan ok", res)
}

// POST /api/v1/ssh/{id}/hostkey/trust
func (h KnownHostsHandler) Trust(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if id <= 0 {
		response.Fail(w, http.StatusBadRequest, "invalid_id", "Invalid id")
		return
	}

	meta, _, err := h.ssh.GetDecrypted(r.Context(), auth.UserID, id)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Fail(w, http.StatusNotFound, "not_found", "Connection not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		return
	}

	// yeniden scan edip ekliyoruz (kopya/yanlış key riskini azaltır)
	scan, err := service.ScanHostKey(r.Context(), meta.Host, meta.Port)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, "scan_failed", err.Error())
		return
	}

	if err := h.kh.EnsureFile(); err != nil {
		response.Fail(w, http.StatusInternalServerError, "server_error", "Cannot prepare known_hosts")
		return
	}
	if err := h.kh.AppendLine(scan.KnownHostsLine); err != nil {
		response.Fail(w, http.StatusInternalServerError, "server_error", "Cannot write known_hosts")
		return
	}

	response.OK(w, "trusted", scan)
}

// DELETE /api/v1/ssh/known-hosts/{token}
func (h KnownHostsHandler) DeleteByToken(w http.ResponseWriter, r *http.Request) {
	_, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	token := chi.URLParam(r, "token")
	if token == "" {
		response.Fail(w, http.StatusBadRequest, "invalid_token", "Invalid host token")
		return
	}

	if err := h.kh.RemoveByHostToken(token); err != nil {
		response.Fail(w, http.StatusInternalServerError, "server_error", "Cannot update known_hosts")
		return
	}
	response.OK(w, "deleted", map[string]any{"token": token})
}
