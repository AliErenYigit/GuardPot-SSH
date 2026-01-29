package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type SecretBox struct {
	key []byte // 32 bytes
}

func NewSecretBox(keyBase64 string) (SecretBox, error) {
	b, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return SecretBox{}, fmt.Errorf("decode SSH_CRED_KEY_BASE64 failed: %w", err)
	}
	if len(b) != 32 {
		return SecretBox{}, fmt.Errorf("SSH_CRED_KEY_BASE64 must decode to 32 bytes (got %d)", len(b))
	}
	return SecretBox{key: b}, nil
}

func (s SecretBox) EncryptToBase64(plain string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("aes cipher failed: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm failed: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce failed: %w", err)
	}

	ct := gcm.Seal(nonce, nonce, []byte(plain), nil) // prepend nonce
	return base64.StdEncoding.EncodeToString(ct), nil
}

func (s SecretBox) DecryptFromBase64(enc string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("aes cipher failed: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm failed: %w", err)
	}

	ns := gcm.NonceSize()
	if len(raw) < ns {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ct := raw[:ns], raw[ns:]
	pt, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt failed: %w", err)
	}
	return string(pt), nil
}
