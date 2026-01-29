package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID int64
	Email  string
	Exp    time.Time
}

// Authorization: Bearer <token> ile gelen JWT'yi doğrular ve claim'leri döner.
func (t TokenManager) VerifyAccessToken(tokenStr string) (AuthClaims, error) {
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Algoritma kontrolü (alg swap saldırılarına karşı)
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return t.secret, nil
	})
	if err != nil || !parsed.Valid {
		return AuthClaims{}, fmt.Errorf("invalid token")
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return AuthClaims{}, fmt.Errorf("invalid claims")
	}

	// sub: string userId
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return AuthClaims{}, fmt.Errorf("missing sub")
	}
	uid, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return AuthClaims{}, fmt.Errorf("invalid sub")
	}

	email, _ := claims["email"].(string)

	// exp: unix
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return AuthClaims{}, fmt.Errorf("missing exp")
	}
	exp := time.Unix(int64(expFloat), 0).UTC()

	return AuthClaims{
		UserID: uid,
		Email:  email,
		Exp:    exp,
	}, nil
}
