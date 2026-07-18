package fetcher

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

func FetchSubscription(subURL string) ([]model.Server, error) {
    resp, err := http.Get(subURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read body: %w", err)
    }

    decoded, err := base64.StdEncoding.DecodeString(string(body))
    if err != nil {
        decoded = body
    }

    lines := strings.Split(string(decoded), "\n")
    var servers []model.Server

    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" {
            continue
        }
        server, err := parseURI(line)
        if err == nil {
            servers = append(servers, server)
        }
    }
    return servers, nil
}

func parseURI(uri string) (model.Server, error) {
    if strings.HasPrefix(uri, "vless://") {
        return parseVless(uri)
    } else if strings.HasPrefix(uri, "vmess://") {
        return parseVmess(uri)
    }
    return model.Server{}, fmt.Errorf("unsupported")
}

func parseVless(uri string) (model.Server, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return model.Server{}, err
    }
    host := u.Hostname()
    port := 0
    fmt.Sscanf(u.Port(), "%d", &port)
    name := u.Query().Get("name")
    if name == "" {
        name = host
    }
    return model.Server{
        ID:       u.Fragment,
        Name:     name,
        Address:  host,
        Port:     port,
        Protocol: "vless",
        URI:      uri,
    }, nil
}

func parseVmess(uri string) (model.Server, error) {
    rawJSON := strings.TrimPrefix(uri, "vmess://")
    decoded, err := base64.StdEncoding.DecodeString(rawJSON)
    if err != nil {
        decoded, err = base64.URLEncoding.DecodeString(rawJSON)
        if err != nil {
            return model.Server{}, fmt.Errorf("b64 decode failed: %w", err)
        }
    }
    var data struct {
        PS   string `json:"ps"`
        Add  string `json:"add"`
        Port string `json:"port"`
    }
    if err := json.Unmarshal(decoded, &data); err != nil {
        return model.Server{}, fmt.Errorf("json unmarshal failed: %w", err)
    }
    port := 0
    fmt.Sscanf(data.Port, "%d", &port)
    return model.Server{
        ID:       data.PS,
        Name:     data.PS,
        Address:  data.Add,
        Port:     port,
        Protocol: "vmess",
        URI:      uri,
    }, nil
}
