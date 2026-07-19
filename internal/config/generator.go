package config

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "strconv"
    "strings"
)

type SingBoxConfig struct {
    Log          map[string]interface{} `json:"log,omitempty"`
    Experimental map[string]interface{} `json:"experimental,omitempty"`
    DNS          map[string]interface{} `json:"dns,omitempty"`
    Inbounds     []map[string]interface{} `json:"inbounds"`
    Outbounds    []map[string]interface{} `json:"outbounds"`
}

func GenerateSingBoxConfig(uri string) (string, string, error) {
    if strings.HasPrefix(uri, "vless://") {
        return generateVlessConfig(uri)
    } else if strings.HasPrefix(uri, "trojan://") {
        return generateTrojanConfig(uri)
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

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return "", "", fmt.Errorf("invalid port format: %w", err)
    }

    outbound := map[string]interface{}{
        "type":        "vless",
        "server":      host,
        "server_port": port,
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
        // FIX BUG-014: Support allowInsecure
        if u.Query().Get("allowInsecure") == "1" {
            tlsOpts["insecure"] = true
        }
        outbound["tls"] = tlsOpts
    }

    transportType := u.Query().Get("type")
    if transportType == "ws" {
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

    // FIX BUG-046: Support UDP-over-TCP
    if u.Query().Get("uot") == "1" {
        outbound["udp_over_tcp"] = true
    }

    return buildConfig(outbound)
}

func generateTrojanConfig(uri string) (string, string, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return "", "", err
    }

    password := u.User.Username()
    host := u.Hostname()
    portStr := u.Port()

    if password == "" || host == "" || portStr == "" {
        return "", "", fmt.Errorf("invalid trojan uri: missing password, host or port")
    }

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return "", "", fmt.Errorf("invalid port format: %w", err)
    }

    outbound := map[string]interface{}{
        "type":        "trojan",
        "server":      host,
        "server_port": port,
        "password":    password,
    }

    sni := u.Query().Get("sni")
    if sni == "" {
        sni = host
    }
    outbound["tls"] = map[string]interface{}{
        "enabled":     true,
        "server_name": sni,
        "insecure":    u.Query().Get("allowInsecure") == "1",
    }

    return buildConfig(outbound)
}

func buildConfig(outbound map[string]interface{}) (string, string, error) {
    clashSecret := generateRandomString(32)
    proxyUser := generateRandomString(16)
    proxyPass := generateRandomString(32)

    config := SingBoxConfig{
        // FIX BUG-034: Set log level to warn to prevent credential leaks
        Log: map[string]interface{}{
            "level": "warn",
            "timestamp": true,
        },
        Experimental: map[string]interface{}{
            "clash_api": map[string]interface{}{
                "external_controller": "127.0.0.1:9090",
                "secret":              clashSecret,
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
                "authentication": []map[string]interface{}{
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
        os.Remove(tmpFile.Name())
        return "", "", err
    }

    if _, err := tmpFile.Write(jsonData); err != nil {
        tmpFile.Close()
        os.Remove(tmpFile.Name())
        return "", "", err
    }
    tmpFile.Close()

    return tmpFile.Name(), clashSecret, nil
}

func GenerateID(uri string) string {
    hash := sha256.Sum256([]byte(uri))
    return hex.EncodeToString(hash[:])
}

func generateRandomString(length int) string {
    b := make([]byte, length)
    rand.Read(b)
    return hex.EncodeToString(b)
}
