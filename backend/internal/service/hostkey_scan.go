package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

type HostKeyScanResult struct {
	HostToken   string `json:"hostToken"`   // example: 192.168.1.50 or [192.168.1.50]:2222
	KeyType     string `json:"keyType"`     // ssh-ed25519 / ssh-rsa ...
	Fingerprint string `json:"fingerprint"` // SHA256:....
	KnownHostsLine string `json:"knownHostsLine"`
}

func HostToken(host string, port int) string {
	if port == 22 {
		return host
	}
	return fmt.Sprintf("[%s]:%d", host, port)
}

func ScanHostKey(ctx context.Context, host string, port int) (HostKeyScanResult, error) {
	addr := SSHAddr(host, port)

	var captured ssh.PublicKey
	cb := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// Handshake sırasında çağrılır
		captured = key
		return nil
	}

	cfg := &ssh.ClientConfig{
		User:            "noop",
		Auth:            []ssh.AuthMethod{}, // auth yok -> handshake sonrası auth fail olabilir, sorun değil
		HostKeyCallback: cb,
		Timeout:         0,
	}

	d := &net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return HostKeyScanResult{}, fmt.Errorf("dial failed: %w", err)
	}
	defer conn.Close()

	// Handshake başlar, host key callback burada tetiklenir
	_, _, _, err = ssh.NewClientConn(conn, addr, cfg)
	_ = err // auth fail olabilir; captured alındıysa yeterli

	if captured == nil {
		return HostKeyScanResult{}, fmt.Errorf("host key could not be captured")
	}

	token := HostToken(host, port)
	line := fmt.Sprintf("%s %s %s", token, captured.Type(), base64.StdEncoding.EncodeToString(captured.Marshal()))

	return HostKeyScanResult{
		HostToken:       token,
		KeyType:         captured.Type(),
		Fingerprint:     ssh.FingerprintSHA256(captured),
		KnownHostsLine:  line,
	}, nil
}
