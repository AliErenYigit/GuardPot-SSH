package middleware

import (
	"context"

	"backend/internal/service"
)

func withAuth(ctx context.Context, claims service.AuthClaims) context.Context {
	return context.WithValue(ctx, AuthKey, claims)
}

func GetAuth(ctx context.Context) (service.AuthClaims, bool) {
	v := ctx.Value(AuthKey)
	if v == nil {
		return service.AuthClaims{}, false
	}
	claims, ok := v.(service.AuthClaims)
	return claims, ok
}
