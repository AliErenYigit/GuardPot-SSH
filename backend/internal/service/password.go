package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordManager struct {
	cost int
}

func NewPasswordManager(cost int) PasswordManager {
	return PasswordManager{cost: cost}
}

func (p PasswordManager) Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), p.cost)
	if err != nil {
		return "", fmt.Errorf("hash password failed: %w", err)
	}
	return string(b), nil
}

func (p PasswordManager) Compare(hash string, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
