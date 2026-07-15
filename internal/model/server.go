package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type Server struct {
	Name     string `json:"name"`
	Config   string `json:"config"`
	Protocol string `json:"protocol"`
	PingMS   int64  `json:"ping_ms"`
	Favorite bool   `json:"favorite"`
	Group    string `json:"group"`
}

func (s Server) Validate() error {
	if strings.TrimSpace(s.Config) == "" {
		return fmt.Errorf("config is empty")
	}
	return nil
}

func ParseServerConfig(config string, group string) Server {
	cfg := strings.TrimSpace(config)
	server := Server{
		Name:     "Custom Server",
		Config:   cfg,
		Protocol: "Unknown",
		Group:    group,
		PingMS:   -1,
	}

	switch {
	case strings.HasPrefix(cfg, "vless://"):
		server.Protocol = "VLESS"
		server.Name = remarkFromURI(cfg, "VLESS Server")
	case strings.HasPrefix(cfg, "vmess://"):
		server.Protocol = "VMess"
		server.Name = remarkFromURI(cfg, "VMess Server")
	case strings.HasPrefix(cfg, "trojan://"):
		server.Protocol = "Trojan"
		server.Name = remarkFromURI(cfg, "Trojan Server")
	case strings.HasPrefix(cfg, "{"):
		var body map[string]any
		if err := json.Unmarshal([]byte(cfg), &body); err == nil {
			if outs, ok := body["outbounds"].([]any); ok && len(outs) > 0 {
				if first, ok := outs[0].(map[string]any); ok {
					if proto, ok := first["protocol"].(string); ok && proto != "" {
						server.Protocol = strings.ToUpper(proto[:1]) + proto[1:]
						server.Name = server.Protocol + " Server"
					}
				}
			}
		}
	}

	if strings.TrimSpace(server.Name) == "" || server.Name == "Custom Server" {
		server.Name = server.Protocol + " Server"
	}
	return server
}

func remarkFromURI(raw string, fallback string) string {
	if idx := strings.LastIndex(raw, "#"); idx >= 0 && idx < len(raw)-1 {
		if decoded, err := url.QueryUnescape(raw[idx+1:]); err == nil && strings.TrimSpace(decoded) != "" {
			return decoded
		}
	}
	return fallback
}
