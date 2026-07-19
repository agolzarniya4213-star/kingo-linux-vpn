package config

import (
    "os"
    "path/filepath"

    "gopkg.in/yaml.v3"
)

type AppConfig struct {
    IPC struct {
        SocketPath string `yaml:"socket_path"`
        MaxConns   int    `yaml:"max_connections"`
    } `yaml:"ipc"`
    Network struct {
        ClashPort  int      `yaml:"clash_api_port"`
        ProxyPort  int      `yaml:"proxy_port"`
        DNSServers []string `yaml:"dns_servers"`
    } `yaml:"network"`
}

func LoadConfig() *AppConfig {
    cfg := &AppConfig{}
    
    // Defaults
    cfg.IPC.SocketPath = "/run/kingo-vpn/kingo-vpn.sock"
    cfg.IPC.MaxConns = 10
    cfg.Network.ClashPort = 9090
    cfg.Network.ProxyPort = 2080
    cfg.Network.DNSServers = []string{"https://1.1.1.1/dns-query", "8.8.8.8"}

    // Try loading from file
    configPath := "/etc/kingo-vpn/config.yaml"
    if os.Geteuid() != 0 {
        if home, err := os.UserHomeDir(); err == nil {
            configPath = filepath.Join(home, ".config", "kingo-vpn", "config.yaml")
        }
    }

    data, err := os.ReadFile(configPath)
    if err == nil {
        yaml.Unmarshal(data, cfg)
    }
    
    return cfg
}
