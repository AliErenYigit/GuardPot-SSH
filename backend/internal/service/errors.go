package service

import "errors"

var (
	ErrEmailAlreadyInUse   = errors.New("email already in use")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrInvalidPassword     = errors.New("invalid password")
)
