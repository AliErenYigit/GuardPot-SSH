package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"backend/internal/config"
	"backend/internal/http/handler"

	"backend/internal/http/response"

	"backend/internal/repository/sqlite"
	"backend/internal/service"

	httpmw "backend/internal/http/middleware"
)

func New(db *sql.DB, cfg config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Get("/api/v1/health", func(w http.ResponseWriter, _ *http.Request) {
		response.OK(w, "ok", map[string]string{"status": "up"})
	})

	// Wiring (repo -> service -> handler)
	userRepo := sqlite.NewUserSQLiteRepo(db)
	pass := service.NewPasswordManager(cfg.BcryptCost)
	tokens := service.NewTokenManager(cfg.JWTSecret, cfg.JWTExpiresMin)
	authSvc := service.NewAuthService(userRepo, pass, tokens)
	authH := handler.NewAuthHandler(authSvc)

	// SSH wiring
	sshRepo := sqlite.NewSSHSQLiteRepo(db)
	auditRepo := sqlite.NewSSHAuditSQLiteRepo(db)

	box, err := service.NewSecretBox(cfg.SSHCredKeyBase64)
	if err != nil {
		panic(err)
	}

	sshSvc := service.NewSSHService(cfg, sshRepo, box)
	sshH := handler.NewSSHHandler(sshSvc)

	limits := httpmw.NewWSLimits(cfg.SSHWSMaxConn, cfg.SSHWSMaxConnPerUser)
	sshWS := handler.NewSSHWSHandler(sshSvc, auditRepo, limits, cfg.SSHKnownHostsPath)

	kh := service.NewKnownHostsManager(cfg.SSHKnownHostsPath)
	khH := handler.NewKnownHostsHandler(sshSvc, kh)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authH.Register)
		r.Post("/login", authH.Login)
	})
	meH := handler.NewMeHandler()

	// Protected endpoints
	r.Group(func(r chi.Router) {
		r.Use(httpmw.JWTAuth(tokens))

		// auth test
		r.Get("/api/v1/me", meH.Me)

		// ssh connections CRUD
		r.Route("/api/v1/ssh", func(r chi.Router) {
			r.Post("/", sshH.Create)
			r.Get("/", sshH.List)
			r.Delete("/{id}", sshH.Delete)

			// live terminal websocket
			r.Get("/{id}/ws", sshWS.Connect)

			// known_hosts management
			r.Get("/known-hosts", khH.List)
			r.Post("/{id}/hostkey/scan", khH.Scan)
			r.Post("/{id}/hostkey/trust", khH.Trust)
			r.Delete("/known-hosts/{token}", khH.DeleteByToken)

		})
	})

	return r
}
