package middleware

import (
	"net/http"
	"strings"

	"backend/internal/http/response"
	"backend/internal/service"
)

func JWTAuth(tokens service.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				response.Fail(w, http.StatusUnauthorized, "missing_token", "Authorization header is required")
				return
			}

			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
				response.Fail(w, http.StatusUnauthorized, "invalid_token", "Authorization header must be Bearer <token>")
				return
			}

			claims, err := tokens.VerifyAccessToken(parts[1])
			if err != nil {
				response.Fail(w, http.StatusUnauthorized, "invalid_token", "Token is invalid or expired")
				return
			}

			ctx := r.Context()
			ctx = withAuth(ctx, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
