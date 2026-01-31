package middleware

import (
	"net/http"
)

func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Origin varsa CORS header set et
			if origin != "" {
				// Dev mod: allowedOrigin="*" ise gelen origin'i aynen yansıt
				if allowedOrigin == "*" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else if origin == allowedOrigin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.Header().Set("Vary", "Origin")

				// Eğer cookie kullanmıyorsan gerek yok, ama ileride lazım olabilir:
				// w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
