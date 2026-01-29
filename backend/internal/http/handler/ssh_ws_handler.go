package handler

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"

	"backend/internal/domain"
	"backend/internal/http/middleware"
	"backend/internal/http/response"
	"backend/internal/repository"
	"backend/internal/service"
	
)

type SSHWSHandler struct {
	ssh    service.SSHService
	audit  repository.SSHAuditRepository
	limits *middleware.WSLimits
	cfg    struct {
		KnownHostsPath string
		Timeout        time.Duration
	}
}

func NewSSHWSHandler(sshSvc service.SSHService, audit repository.SSHAuditRepository, limits *middleware.WSLimits, knownHostsPath string) SSHWSHandler {
	return SSHWSHandler{
		ssh:    sshSvc,
		audit:  audit,
		limits: limits,
		cfg: struct {
			KnownHostsPath string
			Timeout        time.Duration
		}{
			KnownHostsPath: knownHostsPath,
			Timeout:        sshSvc.ConnectTimeout(),
		},
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Prod'da origin kontrolü ekleyeceğiz.
		return true
	},
}

type wsMsg struct {
	Type string `json:"type"` // "data" | "resize"
	Data string `json:"data,omitempty"`

	Cols int `json:"cols,omitempty"`
	Rows int `json:"rows,omitempty"`
}

func (h SSHWSHandler) Connect(w http.ResponseWriter, r *http.Request) {
	auth, ok := middleware.GetAuth(r.Context())
	if !ok {
		response.Fail(w, http.StatusUnauthorized, "unauthorized", "Not authenticated")
		return
	}

	// WS limits
	if h.limits != nil {
		if !h.limits.TryAcquire(auth.UserID) {
			response.Fail(w, http.StatusTooManyRequests, "too_many_connections", "Too many active terminal sessions")
			return
		}
		defer h.limits.Release(auth.UserID)
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(w, http.StatusBadRequest, "invalid_id", "Invalid id")
		return
	}

	connMeta, secret, err := h.ssh.GetDecrypted(r.Context(), auth.UserID, id)
	if err != nil {
		if err == repository.ErrNotFound {
			response.Fail(w, http.StatusNotFound, "not_found", "Connection not found")
			return
		}
		response.Fail(w, http.StatusInternalServerError, "server_error", "Unexpected error")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	remoteIP := clientIP(r)

	// known_hosts callback
	hostKeyCB, err := service.KnownHostsCallback(h.cfg.KnownHostsPath)
	if err != nil {
		_ = h.audit.Log(r.Context(), auth.UserID, id, remoteIP, "connect_fail", err.Error())
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Host key verification not ready: "+err.Error()+"\n"))
		_ = ws.WriteMessage(websocket.TextMessage, []byte("Add server host key to known_hosts and retry.\n"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	client, session, err := dialSSH(ctx, connMeta, secret, hostKeyCB)
	if err != nil {
		_ = h.audit.Log(r.Context(), auth.UserID, id, remoteIP, "connect_fail", err.Error())
		_ = ws.WriteMessage(websocket.TextMessage, []byte("SSH connect failed: "+err.Error()+"\n"))
		return
	}
	defer client.Close()
	defer session.Close()

	_ = h.audit.Log(r.Context(), auth.UserID, id, remoteIP, "connect_ok", "connected")

	// PTY + shell
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	// default size
	rows, cols := 30, 120
	if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		_ = h.audit.Log(r.Context(), auth.UserID, id, remoteIP, "connect_fail", "pty: "+err.Error())
		_ = ws.WriteMessage(websocket.TextMessage, []byte("PTY request failed: "+err.Error()+"\n"))
		return
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("stdin pipe failed: "+err.Error()+"\n"))
		return
	}
	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("stdout pipe failed: "+err.Error()+"\n"))
		return
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("stderr pipe failed: "+err.Error()+"\n"))
		return
	}

	if err := session.Shell(); err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("shell failed: "+err.Error()+"\n"))
		return
	}

	// output -> WS
	done := make(chan struct{})
	go func() {
		defer close(done)
		pipeToWS(ws, stdout)
	}()
	go func() { pipeToWS(ws, stderr) }()

	// input loop
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if len(msg) == 0 {
			continue
		}

		// JSON protocol?
		if looksLikeJSON(msg) {
			var m wsMsg
			if json.Unmarshal(msg, &m) == nil {
				switch strings.ToLower(m.Type) {
				case "resize":
					if m.Cols > 0 && m.Rows > 0 {
						_ = session.WindowChange(m.Rows, m.Cols)
					}
					continue
				case "data":
					if m.Data != "" {
						_, _ = io.WriteString(stdin, m.Data)
					}
					continue
				}
			}
		}

		// fallback: raw text bytes
		_, _ = stdin.Write(msg)
	}

	_ = h.audit.Log(r.Context(), auth.UserID, id, remoteIP, "disconnect", "closed")
	_ = session.Close()
	<-done
}

func dialSSH(ctx context.Context, meta domain.SSHConnection, secret string, hostKeyCB ssh.HostKeyCallback) (*ssh.Client, *ssh.Session, error) {
	var authMethods []ssh.AuthMethod

	switch meta.AuthType {
	case domain.SSHAuthPassword:
		authMethods = append(authMethods, ssh.Password(secret))
	case domain.SSHAuthPrivateKey:
		signer, err := ssh.ParsePrivateKey([]byte(secret))
		if err != nil {
			return nil, nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	default:
		return nil, nil, ssh.ErrNoAuth
	}

	cfg := &ssh.ClientConfig{
		User:            meta.Username,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCB,
		Timeout:         0,
	}

	addr := service.SSHAddr(meta.Host, meta.Port)

	nd := &net.Dialer{}
	c, err := nd.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	cc, chans, reqs, err := ssh.NewClientConn(c, addr, cfg)
	if err != nil {
		_ = c.Close()
		return nil, nil, err
	}
	client := ssh.NewClient(cc, chans, reqs)

	sess, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}
	return client, sess, nil
}

func pipeToWS(ws *websocket.Conn, r io.Reader) {
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_ = ws.WriteMessage(websocket.TextMessage, buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func clientIP(r *http.Request) string {
	// reverse proxy varsa RealIP middleware daha iyi; şimdilik basit
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		parts := strings.Split(ip, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func looksLikeJSON(b []byte) bool {
	s := strings.TrimSpace(string(b))
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}
