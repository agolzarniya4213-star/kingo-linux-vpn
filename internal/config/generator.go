package config

import (
    "crypto/md5"
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "strconv"
    "strings"
)

type SingBoxConfig struct {
    Experimental map[string]interface{} `json:"experimental,omitempty"`
    DNS          map[string]interface{} `json:"dns,omitempty"`
    Inbounds     []map[string]interface{} `json:"inbounds"`
    Outbounds    []map[string]interface{} `json:"outbounds"`
}

func GenerateSingBoxConfig(uri string) (string, string, error) {
    if strings.HasPrefix(uri, "vless://") {
        return generateVlessConfig(uri)
    }
    return "", "", fmt.Errorf("unsupported protocol for config generation")
}

func generateVlessConfig(uri string) (string, string, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return "", "", err
    }

    uuid := u.User.Username()
    host := u.Hostname()
    portStr := u.Port()

    if uuid == "" || host == "" || portStr == "" {
        return "", "", fmt.Errorf("invalid vless uri: missing uuid, host or port")
    }

    // FIX BUG-041: Parse port to integer
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return "", "", fmt.Errorf("invalid port format: %w", err)
    }

    outbound := map[string]interface{}{
        "type":        "vless",
        "server":      host,
        "server_port": port, // Now integer
        "uuid":        uuid,
    }

    security := u.Query().Get("security")
    if security == "tls" {
        sni := u.Query().Get("sni")
        if sni == "" {
            sni = host
        }
        tlsOpts := map[string]interface{}{
            "enabled":     true,
            "server_name": sni,
        }
        if fp := u.Query().Get("fp"); fp != "" {
            tlsOpts["utls"] = map[string]interface{}{"enabled": true, "fingerprint": fp}
        }
        outbound["tls"] = tlsOpts
    }

    transportType := u.Query().Get("type")
    if transportType == "ws" {
        // FIX BUG-044: Fix dead code and properly assign headers
        wsOpts := map[string]interface{}{
            "path": u.Query().Get("path"),
        }
        if hostHeader := u.Query().Get("host"); hostHeader != "" {
            wsOpts["headers"] = map[string]interface{}{"Host": hostHeader}
        }
        outbound["transport"] = wsOpts
    } else if transportType == "grpc" {
        outbound["transport"] = map[string]interface{}{
            "type":         "grpc",
            "service_name": u.Query().Get("serviceName"),
        }
    }

    // FIX BUG-003 & 004: Generate random secrets for Clash API and Local Proxy
    clashSecret := generateRandomString(32)
    proxyUser := generateRandomString(16)
    proxyPass := generateRandomString(32)

    config := SingBoxConfig{
        Experimental: map[string]interface{}{
            "clash_api": map[string]interface{}{
                "external_controller": "127.0.0.1:9090",
                "secret":              clashSecret, // Secured!
            },
        },
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
                "authentication": []map[string]interface{}{ // Secured!
                    {"username": proxyUser, "password": proxyPass},
                },
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
        return "", "", err
    }

    tmpFile, err := os.CreateTemp("", "kingo-config-*.json")
    if err != nil {
        return "", "", err
    }
    
    if err := os.Chmod(tmpFile.Name(), 0600); err != nil {
        tmpFile.Close()
        return "", "", err
    }

    if _, err := tmpFile.Write(jsonData); err != nil {
        tmpFile.Close()
        return "", "", err
    }
    tmpFile.Close()

    return tmpFile.Name(), clashSecret, nil
}

func GenerateID(uri string) string {
    hash := md5.Sum([]byte(uri))
    return hex.EncodeToString(hash[:])
}

func generateRandomString(length int) string {
    b := make([]byte, length)
    rand.Read(b)
    return hex.EncodeToString(b)
}
