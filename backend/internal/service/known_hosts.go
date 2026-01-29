package service

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func KnownHostsCallback(path string) (ssh.HostKeyCallback, error) {
	// Dosya yoksa oluşturmayalım; yoksa açıkça hata verelim.
	// (TOFU eklemek istersek sonraki adımda endpoint ile ekleriz.)
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("known_hosts file not found at %s", path)
	}
	cb, err := knownhosts.New(path)
	if err != nil {
		return nil, fmt.Errorf("known_hosts parse failed: %w", err)
	}
	return cb, nil
}
