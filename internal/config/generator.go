package config

import (
    "encoding/json"
    "fmt"
    "net/url"
    "os"
    "strings"
)

// GenerateSingBoxConfig یک لینک VLESS را گرفته و یک فایل کانفیگ Sing-box می‌سازد
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

    // ساختار مینیمال و معتبر Sing-box
    config := map[string]interface{}{
        "inbounds": []map[string]interface{}{
            {
                "type":        "mixed",
                "listen":      "127.0.0.1",
                "listen_port": 2080,
            },
        },
        "outbounds": []map[string]interface{}{
            {
                "type":        "vless",
                "server":      host,
                "server_port": port,
                "uuid":        uuid,
            },
        },
    }

    // افزودن TLS در صورت وجود پارامتر security=tls
    if u.Query().Get("security") == "tls" {
        sni := u.Query().Get("sni")
        if sni == "" {
            sni = host
        }
        config["outbounds"].([]map[string]interface{})[0]["tls"] = map[string]interface{}{
            "enabled":     true,
            "server_name": sni,
        }
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
