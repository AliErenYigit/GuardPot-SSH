package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type KnownHostsManager struct {
	path string
	mu   sync.Mutex
}

func NewKnownHostsManager(path string) KnownHostsManager {
	return KnownHostsManager{path: path}
}

func (m KnownHostsManager) EnsureFile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, err := os.Stat(m.path); err == nil {
		return nil
	}

	// klasör yoksa oluştur
	dir := "."
	if i := strings.LastIndex(m.path, "/"); i > 0 {
		dir = m.path[:i]
	}
	_ = os.MkdirAll(dir, 0o755)

	f, err := os.OpenFile(m.path, os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	return f.Close()
}

func (m KnownHostsManager) ListLines() ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	f, err := os.Open(m.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, sc.Err()
}

// Aynı satır zaten varsa tekrar eklemez (idempotent)
func (m KnownHostsManager) AppendLine(line string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// dosya yoksa oluştur
	if err := m.EnsureFile(); err != nil {
		return err
	}

	// zaten var mı?
	existing, err := os.ReadFile(m.path)
	if err == nil && strings.Contains(string(existing), line) {
		return nil
	}

	f, err := os.OpenFile(m.path, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintln(f, line)
	return err
}

// host token ile başlayan satırları siler (ör: "192.168.1.50 " veya "[host]:2222 ")
func (m KnownHostsManager) RemoveByHostToken(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, err := os.ReadFile(m.path)
	if err != nil {
		return err
	}
	var out []string
	for _, ln := range strings.Split(string(b), "\n") {
		l := strings.TrimSpace(ln)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, token+" ") {
			continue
		}
		out = append(out, l)
	}
	return os.WriteFile(m.path, []byte(strings.Join(out, "\n")+"\n"), 0o600)
}
