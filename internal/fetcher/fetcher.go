package fetcher

import (
    "context"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "log/slog"
    "net"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/config"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

var httpClient = &http.Client{
    Timeout: 15 * time.Second,
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse // Prevent SSRF via redirects
    },
}

// isPublicIP checks if an IP is public (not loopback, private, or link-local)
func isPublicIP(ip net.IP) bool {
    if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
        return false
    }
    if ip4 := ip.To4(); ip4 != nil {
        if ip4[0] == 10 || (ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || (ip4[0] == 192 && ip4[1] == 168) {
            return false
        }
    }
    return true
}

func FetchSubscription(subURL string) ([]model.Server, error) {
    parsedURL, err := url.Parse(subURL)
    if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
        return nil, fmt.Errorf("invalid URL scheme: only http/https allowed")
    }

    // FIX BUG-005: SSRF Prevention
    ips, err := net.LookupIP(parsedURL.Hostname())
    if err != nil {
        return nil, fmt.Errorf("failed to resolve hostname: %w", err)
    }
    for _, ip := range ips {
        if !isPublicIP(ip) {
            return nil, fmt.Errorf("internal IP addresses are blocked to prevent SSRF")
        }
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    req, err := http.NewRequestWithContext(ctx, "GET", subURL, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    resp, err := httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("subscription server returned error: %s", resp.Status)
    }

    body, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
    if err != nil {
        return nil, fmt.Errorf("failed to read body: %w", err)
    }

    decoded, err := base64.StdEncoding.DecodeString(string(body))
    if err != nil {
        slog.Warn("Base64 decode failed, falling back to plain text")
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
    } else if strings.HasPrefix(uri, "trojan://") {
        return parseTrojan(uri)
    }
    return model.Server{}, fmt.Errorf("unsupported")
}

func parseVless(uri string) (model.Server, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return model.Server{}, err
    }
    host := u.Hostname()
    portStr := u.Port()
    if host == "" || portStr == "" {
        return model.Server{}, fmt.Errorf("invalid vless uri")
    }
    
    // FIX BUG-049: Use strconv.Atoi with error checking
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return model.Server{}, fmt.Errorf("invalid port format: %w", err)
    }
    
    name := u.Fragment
    if name == "" {
        name = host
    }
    
    return model.Server{
        ID:       config.GenerateID(uri),
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
    port, err := strconv.Atoi(data.Port)
    if err != nil {
        return model.Server{}, fmt.Errorf("invalid port format: %w", err)
    }
    if data.Add == "" || port == 0 {
        return model.Server{}, fmt.Errorf("invalid vmess uri")
    }
    return model.Server{
        ID:       config.GenerateID(uri),
        Name:     data.PS,
        Address:  data.Add,
        Port:     port,
        Protocol: "vmess",
        URI:      uri,
    }, nil
}

// FIX BUG-042: Added Trojan parser
func parseTrojan(uri string) (model.Server, error) {
    u, err := url.Parse(uri)
    if err != nil {
        return model.Server{}, err
    }
    host := u.Hostname()
    portStr := u.Port()
    if host == "" || portStr == "" {
        return model.Server{}, fmt.Errorf("invalid trojan uri")
    }
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return model.Server{}, fmt.Errorf("invalid port format: %w", err)
    }
    name := u.Fragment
    if name == "" {
        name = host
    }
    return model.Server{
        ID:       config.GenerateID(uri),
        Name:     name,
        Address:  host,
        Port:     port,
        Protocol: "trojan",
        URI:      uri,
    }, nil
}
