package config

import (
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "strings"
)

type SingBoxConfig struct {
    Experimental map[string]interface{} `json:"experimental,omitempty"`
    DNS          map[string]interface{} `json:"dns,omitempty"`
    Inbounds     []map[string]interface{} `json:"inbounds"`
    Outbounds    []map[string]interface{} `json:"outbounds"`
}

func GenerateSingBoxConfig(uri string) (string, error) {
    if strings.HasPrefix(uri, "vless://") {
        return generateVlessConfig(uri)
    }
    return "", fmt.Errorf("unsupported protocol for config generation")
}

func generateVlessConfig(uri string) (string, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return "", err
    }

    uuid := u.User.Username()
    host := u.Hostname()
    port := u.Port()

    if uuid == "" || host == "" || port == "" {
        return "", fmt.Errorf("invalid vless uri: missing uuid, host or port")
    }

    outbound := map[string]interface{}{
        "type":        "vless",
        "server":      host,
        "server_port": port,
        "uuid":        uuid,
    }

    if u.Query().Get("security") == "tls" {
        sni := u.Query().Get("sni")
        if sni == "" {
            sni = host
        }
        outbound["tls"] = map[string]interface{}{
            "enabled":     true,
            "server_name": sni,
        }
    }

    config := SingBoxConfig{
        Experimental: map[string]interface{}{
            "clash_api": map[string]interface{}{
                "external_controller": "127.0.0.1:9090",
            },
        },
        // جلوگیری از DNS Leak: تمام درخواست‌های DNS از داخل تونل عبور می‌کنند
        DNS: map[string]interface{}{
            "servers": []map[string]interface{}{
                {"address": "https://1.1.1.1/dns-query", "detour": "proxy"},
                {"address": "8.8.8.8", "detour": "direct"},
            },
            "final": "proxy_dns",
            "rules": []map[string]interface{}{
                {"outbound": "any", "server": "direct"},
            },
        },
        Inbounds: []map[string]interface{}{
            {
                "type":        "mixed",
                "listen":      "127.0.0.1",
                "listen_port": 2080,
            },
        },
        Outbounds: []map[string]interface{}{
            outbound,
            {"type": "direct", "tag": "direct"},
            {"type": "dns", "tag": "proxy_dns"},
        },
    }

    jsonData, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return "", err
    }

    // ایجاد فایل موقت با دسترسی محدود (0600) برای جلوگیری از نشت UUID
    tmpFile, err := os.CreateTemp("", "kingo-config-*.json")
    if err != nil {
        return "", err
    }
    
    // اعمال دسترسی امنیتی
    if err := os.Chmod(tmpFile.Name(), 0600); err != nil {
        tmpFile.Close()
        return "", err
    }

    if _, err := tmpFile.Write(jsonData); err != nil {
        tmpFile.Close()
        return "", err
    }
    tmpFile.Close()

    return tmpFile.Name(), nil
}
