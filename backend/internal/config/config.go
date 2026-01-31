package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBPath string

	JWTSecret     string
	JWTExpiresMin int
	BcryptCost    int

	SSHCredKeyBase64     string
	SSHConnectTimeoutSec int

	SSHKnownHostsPath   string
	SSHWSMaxConn        int
	SSHWSMaxConnPerUser int
	 CORSAllowedOrigin string
}

func Load() (Config, error) {
	// .env yoksa da patlamasın (prod'da env üzerinden gider)
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:  getEnv("APP_ENV", "dev"),
		AppPort: getEnv("APP_PORT", "8080"),
		DBPath:  getEnv("DB_PATH", "./data/app.db"),
	}

	cfg.JWTSecret = getEnv("JWT_SECRET", "")
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("missing required env: JWT_SECRET")
	}

	cfg.JWTExpiresMin = getEnvInt("JWT_EXPIRES_MIN", 60)
	cfg.BcryptCost = getEnvInt("BCRYPT_COST", 12)

	// basit güvenlik/sağlamlık kontrolleri
	if cfg.JWTExpiresMin <= 0 {
		return Config{}, fmt.Errorf("JWT_EXPIRES_MIN must be > 0")
	}
	if cfg.BcryptCost < 10 || cfg.BcryptCost > 15 {
		return Config{}, fmt.Errorf("BCRYPT_COST must be between 10 and 15 (got %d)", cfg.BcryptCost)
	}

	cfg.SSHCredKeyBase64 = getEnv("SSH_CRED_KEY_BASE64", "")
	if cfg.SSHCredKeyBase64 == "" {
		return Config{}, fmt.Errorf("missing required env: SSH_CRED_KEY_BASE64")
	}
	cfg.SSHConnectTimeoutSec = getEnvInt("SSH_CONNECT_TIMEOUT_SEC", 10)
	if cfg.SSHConnectTimeoutSec <= 0 {
		cfg.SSHConnectTimeoutSec = 10
	}
	cfg.CORSAllowedOrigin=getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173")
	

	cfg.SSHKnownHostsPath = getEnv("SSH_KNOWN_HOSTS_PATH", "./data/known_hosts")
	cfg.SSHWSMaxConn = getEnvInt("SSH_WS_MAX_CONN", 20)
	if cfg.SSHWSMaxConn <= 0 {
		cfg.SSHWSMaxConn = 20
	}
	cfg.SSHWSMaxConnPerUser = getEnvInt("SSH_WS_MAX_CONN_PER_USER", 3)
	if cfg.SSHWSMaxConnPerUser <= 0 {
		cfg.SSHWSMaxConnPerUser = 3
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}
