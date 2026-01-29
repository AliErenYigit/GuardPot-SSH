package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret     []byte
	expiresMin int
}

func NewTokenManager(secret string, expiresMin int) TokenManager {
	return TokenManager{
		secret:     []byte(secret),
		expiresMin: expiresMin,
	}
}

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}

func (t TokenManager) GenerateAccessToken(userID int64, email string) (AccessToken, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Duration(t.expiresMin) * time.Minute)

	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", userID),
		"email": email,
		"iat":   now.Unix(),
		"exp":   exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(t.secret)
	if err != nil {
		return AccessToken{}, fmt.Errorf("sign token failed: %w", err)
	}

	return AccessToken{Token: signed, ExpiresAt: exp}, nil
}
