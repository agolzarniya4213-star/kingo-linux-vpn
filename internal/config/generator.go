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
        Inbounds: []map[string]interface{}{
            {
                "type":        "mixed",
                "listen":      "127.0.0.1",
                "listen_port": 2080,
            },
        },
        Outbounds: []map[string]interface{}{outbound},
    }

    jsonData, err := json.MarshalIndent(config, "", "  ")
    if err != nil {
        return "", err
    }

    tmpFile, err := os.CreateTemp("", "kingo-config-*.json")
    if err != nil {
        return "", err
    }
    defer tmpFile.Close()

    if _, err := tmpFile.Write(jsonData); err != nil {
        return "", err
    }

    return tmpFile.Name(), nil
}
