package core

import (
	"encoding/json"
	"fmt"

	"github.com/kingo-linux/vpn/internal/config"
	"github.com/kingo-linux/vpn/internal/model"
)

type xrayConfig struct {
	Log       logSection     `json:"log"`
	Inbounds  []inbound      `json:"inbounds"`
	Outbounds []outbound     `json:"outbounds"`
	Routing   map[string]any `json:"routing,omitempty"`
	DNS       map[string]any `json:"dns,omitempty"`
}

type logSection struct {
	Loglevel string `json:"loglevel"`
}

type inbound struct {
	Listen         string         `json:"listen"`
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Settings       map[string]any `json:"settings,omitempty"`
	Sniffing       map[string]any `json:"sniffing,omitempty"`
	StreamSettings map[string]any `json:"streamSettings,omitempty"`
	Tag            string         `json:"tag,omitempty"`
}

type outbound struct {
	Protocol       string         `json:"protocol"`
	Settings       map[string]any `json:"settings,omitempty"`
	StreamSettings map[string]any `json:"streamSettings,omitempty"`
	Tag            string         `json:"tag,omitempty"`
}

func BuildXrayConfig(server model.Server, settings config.Settings) ([]byte, error) {
	if err := server.Validate(); err != nil {
		return nil, err
	}
	ob, err := buildOutbound(server)
	if err != nil {
		return nil, err
	}

	cfg := xrayConfig{
		Log: logSection{Loglevel: "warning"},
		Inbounds: []inbound{
			{
				Listen:   "127.0.0.1",
				Port:     settings.SOCKSPort,
				Protocol: "socks",
				Settings: map[string]any{
					"udp": true,
				},
				Sniffing: map[string]any{
					"enabled":      true,
					"destOverride": []string{"http", "tls"},
					"routeOnly":    false,
				},
				Tag: "socks-in",
			},
		},
		Outbounds: []outbound{
			ob,
			{
				Protocol: "freedom",
				Tag:      "direct",
			},
			{
				Protocol: "blackhole",
				Tag:      "block",
			},
		},
		Routing: map[string]any{
			"domainStrategy": "AsIs",
			"rules": []map[string]any{
				{
					"type":        "field",
					"inboundTag":  []string{"socks-in"},
					"outboundTag": ob.Tag,
					"network":     "tcp,udp",
				},
			},
		},
		DNS: map[string]any{
			"servers": []any{"1.1.1.1", "8.8.8.8"},
		},
	}
	return json.MarshalIndent(cfg, "", "  ")
}

func buildOutbound(server model.Server) (outbound, error) {
	switch server.Protocol {
	case "VLESS":
		return outbound{
			Protocol: "vless",
			Tag:      "proxy",
			Settings: map[string]any{
				"vnext": []any{
					map[string]any{
						"address": "example.com",
						"port":    443,
						"users": []any{
							map[string]any{
								"id":         "00000000-0000-0000-0000-000000000000",
								"encryption": "none",
							},
						},
					},
				},
			},
			StreamSettings: map[string]any{
				"network":  "tcp",
				"security": "tls",
			},
		}, nil
	case "VMess":
		return outbound{
			Protocol: "vmess",
			Tag:      "proxy",
			Settings: map[string]any{
				"vnext": []any{
					map[string]any{
						"address": "example.com",
						"port":    443,
						"users": []any{
							map[string]any{
								"id":       "00000000-0000-0000-0000-000000000000",
								"alterId":  0,
								"security": "auto",
							},
						},
					},
				},
			},
			StreamSettings: map[string]any{
				"network":  "tcp",
				"security": "tls",
			},
		}, nil
	case "Trojan":
		return outbound{
			Protocol: "trojan",
			Tag:      "proxy",
			Settings: map[string]any{
				"servers": []any{
					map[string]any{
						"address":  "example.com",
						"port":     443,
						"password": "password",
					},
				},
			},
			StreamSettings: map[string]any{
				"network":  "tcp",
				"security": "tls",
			},
		}, nil
	default:
		return outbound{}, fmt.Errorf("unsupported protocol: %s", server.Protocol)
	}
}
