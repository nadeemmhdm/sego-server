package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	ModeLowResource = "low-resource"
	ModeStandard    = "standard"
)

type Config struct {
	ListenAddr        string        `json:"listen_addr"`
	DatabasePath      string        `json:"database_path"`
	SessionHours      int           `json:"session_hours"`
	MaxUploadBytes    int64         `json:"max_upload_bytes"`
	UploadConcurrency int           `json:"upload_concurrency"`
	LowResourceMode   bool          `json:"low_resource_mode"`
	StorageRoots      []StorageRoot `json:"storage_roots"`
}

type StorageRoot struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func Default(base string) Config {
	if base == "" {
		base = defaultBaseDir()
	}
	return Config{
		ListenAddr:        "127.0.0.1:8088",
		DatabasePath:      filepath.Join(base, "sego.db"),
		SessionHours:      24,
		MaxUploadBytes:    10 << 30,
		UploadConcurrency: 2,
		LowResourceMode:   true,
		StorageRoots:      []StorageRoot{{ID: "primary", Name: "Primary", Path: filepath.Join(base, "storage")}},
	}
}

func defaultBaseDir() string {
	if runtime.GOOS == "windows" {
		if v := os.Getenv("ProgramData"); v != "" {
			return filepath.Join(v, "SEGO")
		}
	}
	if runtime.GOOS == "linux" {
		return "/var/lib/sego"
	}
	d, _ := os.UserConfigDir()
	return filepath.Join(d, "sego")
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	if err := normalize(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

func Ensure(path string) (Config, error) {
	c, err := Load(path)
	if err == nil {
		return c, nil
	}
	if !os.IsNotExist(err) {
		return Config{}, err
	}
	c = Default(filepath.Dir(path))
	if err := ensureDirectories(c); err != nil {
		return Config{}, err
	}
	if err := Save(path, c); err != nil {
		return Config{}, err
	}
	return c, nil
}

func Save(path string, c Config) error {
	if err := normalize(&c); err != nil {
		return err
	}
	if err := ensureDirectories(c); err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err == nil {
		return nil
	} else if runtime.GOOS != "windows" {
		_ = os.Remove(tmp)
		return err
	}

	// Windows does not reliably replace an existing destination with os.Rename.
	// Use a rollback file so a failed replacement does not silently destroy the
	// last known-good configuration.
	backup := path + ".bak"
	_ = os.Remove(backup)
	if _, err := os.Stat(path); err == nil {
		if err := os.Rename(path, backup); err != nil {
			_ = os.Remove(tmp)
			return err
		}
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Rename(backup, path)
		return err
	}
	_ = os.Remove(backup)
	return nil
}

func Mode(c Config) string {
	if c.LowResourceMode {
		return ModeLowResource
	}
	return ModeStandard
}

func ApplyMode(c *Config, mode string) error {
	switch mode {
	case ModeLowResource:
		c.LowResourceMode = true
		c.UploadConcurrency = 1
	case ModeStandard:
		c.LowResourceMode = false
		if c.UploadConcurrency < 3 {
			c.UploadConcurrency = 3
		}
	default:
		return fmt.Errorf("unknown mode %q (expected %s or %s)", mode, ModeLowResource, ModeStandard)
	}
	return nil
}

func normalize(c *Config) error {
	if c.ListenAddr == "" || c.DatabasePath == "" || len(c.StorageRoots) == 0 {
		return errors.New("invalid config: listen_addr, database_path and storage_roots are required")
	}
	if c.UploadConcurrency < 1 {
		c.UploadConcurrency = 1
	}
	if c.SessionHours < 1 {
		c.SessionHours = 24
	}
	if c.MaxUploadBytes < 1 {
		c.MaxUploadBytes = 1 << 30
	}
	return nil
}

func ensureDirectories(c Config) error {
	if err := os.MkdirAll(filepath.Dir(c.DatabasePath), 0o700); err != nil {
		return err
	}
	for _, r := range c.StorageRoots {
		if err := os.MkdirAll(r.Path, 0o750); err != nil {
			return err
		}
	}
	return nil
}
