package config

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
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

func GenerateSingBoxConfig(uri string, cfg *AppConfig) (string, string, error) {
    if strings.HasPrefix(uri, "vless://") {
        return generateVlessConfig(uri, cfg)
    } else if strings.HasPrefix(uri, "trojan://") {
        return generateTrojanConfig(uri, cfg)
    } else if strings.HasPrefix(uri, "vmess://") {
        return generateVmessConfig(uri, cfg)
    }
    return "", "", fmt.Errorf("unsupported protocol")
}

func generateVlessConfig(uri string, cfg *AppConfig) (string, string, error) {
    u, err := url.Parse(uri)
    if err != nil { return "", "", err }

    uuid := u.User.Username()
    host := u.Hostname()
    portStr := u.Port()
    if uuid == "" || host == "" || portStr == "" { return "", "", fmt.Errorf("invalid vless uri") }

    port, err := strconv.Atoi(portStr)
    if err != nil { return "", "", fmt.Errorf("invalid port format") }

    outbound := map[string]interface{}{
        "type": "vless", "tag": "proxy",
        "server": host, "server_port": port, "uuid": uuid,
    }

    if flow := u.Query().Get("flow"); flow != "" {
        outbound["flow"] = flow
    }

    security := u.Query().Get("security")
    if security == "tls" {
        sni := u.Query().Get("sni")
        if sni == "" { sni = host }
        tlsOpts := map[string]interface{}{"enabled": true, "server_name": sni}
        if fp := u.Query().Get("fp"); fp != "" {
            tlsOpts["utls"] = map[string]interface{}{"enabled": true, "fingerprint": fp}
        }
        if u.Query().Get("allowInsecure") == "1" { tlsOpts["insecure"] = true }
        outbound["tls"] = tlsOpts
    } else if security == "reality" {
        pbk := u.Query().Get("pbk")
        sid := u.Query().Get("sid")
        if pbk == "" { return "", "", fmt.Errorf("missing pbk for reality") }
        
        sni := u.Query().Get("sni")
        if sni == "" { sni = host }
        
        realityOpts := map[string]interface{}{"enabled": true, "public_key": pbk}
        if sid != "" { realityOpts["short_id"] = sid }
        
        outbound["tls"] = map[string]interface{}{
            "enabled": true, "server_name": sni,
            "reality": realityOpts,
            "utls": map[string]interface{}{"enabled": true, "fingerprint": u.Query().Get("fp")},
        }
    }

    if t := u.Query().Get("type"); t == "ws" {
        wsOpts := map[string]interface{}{"path": u.Query().Get("path")}
        if h := u.Query().Get("host"); h != "" { wsOpts["headers"] = map[string]interface{}{"Host": h} }
        outbound["transport"] = wsOpts
    } else if t == "grpc" {
        outbound["transport"] = map[string]interface{}{"type": "grpc", "service_name": u.Query().Get("serviceName")}
    }

    return buildConfig(outbound, cfg)
}

func generateTrojanConfig(uri string, cfg *AppConfig) (string, string, error) {
    u, err := url.Parse(uri)
    if err != nil { return "", "", err }

    password := u.User.Username()
    host := u.Hostname()
    portStr := u.Port()
    if password == "" || host == "" || portStr == "" { return "", "", fmt.Errorf("invalid trojan uri") }

    port, err := strconv.Atoi(portStr)
    if err != nil { return "", "", fmt.Errorf("invalid port format") }

    outbound := map[string]interface{}{
        "type": "trojan", "tag": "proxy",
        "server": host, "server_port": port, "password": password,
    }
    sni := u.Query().Get("sni")
    if sni == "" { sni = host }
    outbound["tls"] = map[string]interface{}{
        "enabled": true, "server_name": sni, "insecure": u.Query().Get("allowInsecure") == "1",
    }

    return buildConfig(outbound, cfg)
}

func generateVmessConfig(uri string, cfg *AppConfig) (string, string, error) {
    rawJSON := strings.TrimPrefix(uri, "vmess://")
    decoded, err := base64.StdEncoding.DecodeString(rawJSON)
    if err != nil {
        decoded, err = base64.URLEncoding.DecodeString(rawJSON)
        if err != nil { return "", "", fmt.Errorf("b64 decode failed") }
    }

    var data struct {
        PS   string `json:"ps"`
        Add  string `json:"add"`
        Port string `json:"port"`
        ID   string `json:"id"`
        AID  string `json:"aid"`
        Net  string `json:"net"`
        Path string `json:"path"`
        Host string `json:"host"`
        TLS  string `json:"tls"`
        SNI  string `json:"sni"`
    }
    if err := json.Unmarshal(decoded, &data); err != nil { return "", "", fmt.Errorf("json unmarshal failed") }

    port, err := strconv.Atoi(data.Port)
    if err != nil { return "", "", fmt.Errorf("invalid port format") }

    outbound := map[string]interface{}{
        "type": "vmess", "tag": "proxy",
        "server": data.Add, "server_port": port, 
        "uuid": data.ID, "alter_id": data.AID, "security": "auto",
    }

    if data.TLS == "tls" {
        sni := data.SNI
        if sni == "" { sni = data.Add }
        outbound["tls"] = map[string]interface{}{"enabled": true, "server_name": sni, "insecure": false}
    }

    if data.Net == "ws" {
        wsOpts := map[string]interface{}{"path": data.Path}
        if data.Host != "" { wsOpts["headers"] = map[string]interface{}{"Host": data.Host} }
        outbound["transport"] = wsOpts
    } else if data.Net == "grpc" {
        outbound["transport"] = map[string]interface{}{"type": "grpc", "service_name": data.Path}
    }

    return buildConfig(outbound, cfg)
}

func buildConfig(outbound map[string]interface{}, cfg *AppConfig) (string, string, error) {
    clashSecret := generateRandomString(32)

    // FIX: TUN Inbound for System-wide traffic routing
    inbounds := []map[string]interface{}{
        {
            "type": "tun", "tag": "tun-in",
            "interface_name": "kingo0", "auto_route": true, "auto_detect_interface": true,
        },
    }

    // FIX: Explicit DNS tags for sing-box v1.8
    dnsServers := []map[string]interface{}{
        {"tag": "dns-remote", "address": "https://1.1.1.1/dns-query", "detour": "proxy"},
        {"tag": "dns-direct", "address": "8.8.8.8", "detour": "direct"},
    }

    config := SingBoxConfig{
        Log: map[string]interface{}{"level": "warn", "timestamp": true},
        Experimental: map[string]interface{}{
            "clash_api": map[string]interface{}{
                "external_controller": fmt.Sprintf("127.0.0.1:%d", cfg.Network.ClashPort),
                "secret": clashSecret,
            },
        },
        DNS: map[string]interface{}{
            "servers": dnsServers,
            "rules": []map[string]interface{}{{"outbound": "any", "server": "dns-direct"}},
            "final": "dns-remote", 
            "strategy": "ipv4_only",
        },
        Inbounds: inbounds,
        Outbounds: []map[string]interface{}{
            outbound, 
            {"type": "direct", "tag": "direct"}, 
            {"type": "dns", "tag": "dns-out"},
        },
    }

    jsonData, err := json.MarshalIndent(config, "", "  ")
    if err != nil { return "", "", err }

    tmpFile, err := os.CreateTemp("", "kingo-config-*.json")
    if err != nil { return "", "", err }
    
    if err := os.Chmod(tmpFile.Name(), 0600); err != nil {
        tmpFile.Close(); os.Remove(tmpFile.Name()); return "", "", err
    }
    if _, err := tmpFile.Write(jsonData); err != nil {
        tmpFile.Close(); os.Remove(tmpFile.Name()); return "", "", err
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
