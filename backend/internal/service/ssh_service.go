package service

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/domain"
	"backend/internal/repository"
)

type SSHService struct {
	cfg   config.Config
	repo  repository.SSHConnectionRepository
	box   SecretBox
}

func NewSSHService(cfg config.Config, repo repository.SSHConnectionRepository, box SecretBox) SSHService {
	return SSHService{cfg: cfg, repo: repo, box: box}
}

type SSHConnectionPublic struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Username  string    `json:"username"`
	AuthType  string    `json:"authType"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s SSHService) Create(ctx context.Context, userID int64, name, host string, port int, username string, authType domain.SSHAuthType, secretPlain string) (SSHConnectionPublic, error) {
	name = strings.TrimSpace(name)
	host = strings.TrimSpace(host)
	username = strings.TrimSpace(username)

	if name == "" || host == "" || username == "" {
		return SSHConnectionPublic{}, fmt.Errorf("invalid input")
	}
	if port <= 0 || port > 65535 {
		return SSHConnectionPublic{}, fmt.Errorf("invalid port")
	}
	// basit host kontrolü (ip/hostname olabilir)
	if len(host) < 1 || len(host) > 255 {
		return SSHConnectionPublic{}, fmt.Errorf("invalid host")
	}
	if authType != domain.SSHAuthPassword && authType != domain.SSHAuthPrivateKey {
		return SSHConnectionPublic{}, fmt.Errorf("invalid authType")
	}
	if secretPlain == "" {
		return SSHConnectionPublic{}, fmt.Errorf("secret required")
	}

	// küçük bir sanity check: host çözümlenebiliyor mu? (opsiyonel ama pratik)
	_, _ = net.LookupHost(host)

	enc, err := s.box.EncryptToBase64(secretPlain)
	if err != nil {
		return SSHConnectionPublic{}, err
	}

	created, err := s.repo.Create(ctx, domain.SSHConnection{
		UserID:    userID,
		Name:      name,
		Host:      host,
		Port:      port,
		Username:  username,
		AuthType:  authType,
		SecretEnc: enc,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return SSHConnectionPublic{}, err
	}

	return SSHConnectionPublic{
		ID:        created.ID,
		Name:      created.Name,
		Host:      created.Host,
		Port:      created.Port,
		Username:  created.Username,
		AuthType:  string(created.AuthType),
		CreatedAt: created.CreatedAt,
	}, nil
}

func (s SSHService) List(ctx context.Context, userID int64) ([]SSHConnectionPublic, error) {
	items, err := s.repo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]SSHConnectionPublic, 0, len(items))
	for _, c := range items {
		out = append(out, SSHConnectionPublic{
			ID:        c.ID,
			Name:      c.Name,
			Host:      c.Host,
			Port:      c.Port,
			Username:  c.Username,
			AuthType:  string(c.AuthType),
			CreatedAt: c.CreatedAt,
		})
	}
	return out, nil
}

func (s SSHService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s SSHService) GetDecrypted(ctx context.Context, userID, id int64) (domain.SSHConnection, string, error) {
	c, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return domain.SSHConnection{}, "", err
	}
	plain, err := s.box.DecryptFromBase64(c.SecretEnc)
	if err != nil {
		return domain.SSHConnection{}, "", err
	}
	return c, plain, nil
}

// timeout helper
func (s SSHService) ConnectTimeout() time.Duration {
	return time.Duration(s.cfg.SSHConnectTimeoutSec) * time.Second
}

// addr helper
func SSHAddr(host string, port int) string {
	return net.JoinHostPort(host, fmt.Sprintf("%d", port))
}
