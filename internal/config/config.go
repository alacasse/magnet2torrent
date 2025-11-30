package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config captures user-adjustable settings.
type Config struct {
	SaveDir    string `json:"saveDir"`
	LogLevel   string `json:"logLevel"`
	AppName    string `json:"appName"`
	QbUsername string `json:"qbUsername"`
	QbPassword string `json:"qbPassword"`
	QbHost     string `json:"qbHost"`
}

// DefaultConfig returns a config populated with sensible defaults.
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return defaultConfig(home)
}

// GetDefaultConfigPath selects the platform-appropriate config location.
func GetDefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	appdata := os.Getenv("APPDATA")
	return defaultConfigPath(runtime.GOOS, home, appdata)
}

// LoadConfig attempts to read a JSON config; if missing, defaults are returned.
// The returned boolean is true when defaults were used (file missing).
func LoadConfig(path string) (*Config, bool, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path) // #nosec G304 - user-provided path is expected.
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, true, nil
		}
		return nil, false, fmt.Errorf("read config %s: %w", path, err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, false, fmt.Errorf("parse config %s: %w", path, err)
	}

	return cfg, false, nil
}

func defaultConfig(home string) *Config {
	return &Config{
		SaveDir:  defaultSaveDir(home),
		LogLevel: "info",
		AppName:  "magnet2torrent",
	}
}

func defaultConfigPath(goos string, home string, appdata string) string {
	if goos == "windows" {
		base := appdata
		if base == "" {
			base = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(base, "magnet2torrent", "config.json")
	}

	return filepath.Join(home, ".config", "magnet2torrent", "config.json")
}

func defaultSaveDir(home string) string {
	if home == "" {
		return "magnet2torrent-downloads"
	}
	return filepath.Join(home, "Downloads", "magnet2torrent")
}
