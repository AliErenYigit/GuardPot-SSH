package middleware

import "net/http"

// If WebSocket client cannot set Authorization header, allow token via ?token=...
func WSQueryTokenToAuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If header already exists, do nothing
		if r.Header.Get("Authorization") == "" {
			if tok := r.URL.Query().Get("token"); tok != "" {
				r.Header.Set("Authorization", "Bearer "+tok)
			}
		}
		next.ServeHTTP(w, r)
	})
}
