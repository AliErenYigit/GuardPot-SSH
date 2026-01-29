package domain

import "time"

type SSHAuthType string

const (
	SSHAuthPassword   SSHAuthType = "password"
	SSHAuthPrivateKey SSHAuthType = "private_key"
)

type SSHConnection struct {
	ID        int64
	UserID    int64
	Name      string
	Host      string
	Port      int
	Username  string
	AuthType  SSHAuthType
	SecretEnc string // encrypted secret (base64)
	CreatedAt time.Time
}
