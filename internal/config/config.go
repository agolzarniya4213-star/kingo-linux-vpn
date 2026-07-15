package config

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	AppDir          string
	ServersFile     string
	SettingsFile    string
	CacheFile       string
	EngineBinary    string
	SubscriptionURL string
}

type Settings struct {
	SOCKSPort int `json:"socks_port"`
	HTTPPort  int `json:"http_port"`
	DNSPort   int `json:"dns_port"`
}

func Default() Config {
	dir := appDir()
	return Config{
		AppDir:          dir,
		ServersFile:     filepath.Join(dir, "servers.json"),
		SettingsFile:    filepath.Join(dir, "settings.json"),
		CacheFile:       filepath.Join(dir, "cache.json"),
		EngineBinary:    filepath.Join(dir, "bin", "xray"),
		SubscriptionURL: "https://raw.githubusercontent.com/kingowow/Kingo-vpn/refs/heads/main/server/KingoVpn.txt",
	}
}

func DefaultSettings() Settings {
	return Settings{
		SOCKSPort: 10808,
		HTTPPort:  10809,
		DNSPort:   53,
	}
}

func appDir() string {
	if v := strings.TrimSpace(os.Getenv("KINGO_LINUX_VPN_DIR")); v != "" {
		return v
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", "kingo-linux-vpn")
	}
	return ".kingo-linux-vpn"
}
