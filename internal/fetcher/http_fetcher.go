package fetcher

import (
    "bufio"
    "context"
    "fmt"
    "net/http"
    "net/url"
    "strings"

    "kingo-linux-vpn/internal/model"
)

type HttpFetcher struct {
    client *http.Client
}

func NewHttpFetcher() *HttpFetcher {
    return &HttpFetcher{client: &http.Client{}}
}

func (f *HttpFetcher) Fetch(ctx context.Context, subURL string) ([]model.Server, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", subURL, nil)
    if err != nil {
        return nil, fmt.Errorf("create request failed: %w", err)
    }

    resp, err := f.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("fetch failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
    }

    var servers []model.Server
    scanner := bufio.NewScanner(resp.Body)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        server, err := parseLine(line)
        if err != nil {
            continue
        }
        servers = append(servers, server)
    }
    return servers, scanner.Err()
}

func parseLine(line string) (model.Server, error) {
    u, err := url.Parse(line)
    if err != nil {
        return model.Server{}, err
    }

    hostPort := strings.Split(u.Host, ":")
    if len(hostPort) != 2 {
        return model.Server{}, fmt.Errorf("invalid host:port")
    }

    var port int
    _, err = fmt.Sscanf(hostPort[1], "%d", &port)
    if err != nil {
        return model.Server{}, err
    }

    name := u.Fragment
    if name == "" {
        name = u.Host
    }

    return model.Server{
        Name:     name,
        Protocol: strings.ToLower(u.Scheme),
        Address:  hostPort[0],
        Port:     port,
        RawURI:   line,
    }, nil
}
