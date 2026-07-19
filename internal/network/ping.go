package network

import (
    "context"
    "net"
    "strconv"
    "sync"
    "time"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

func TestLatency(ctx context.Context, host string, port int) (int, error) {
    address := net.JoinHostPort(host, strconv.Itoa(port))
    start := time.Now()

    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    d := net.Dialer{}
    conn, err := d.DialContext(ctx, "tcp", address)
    if err != nil {
        return 0, err
    }
    defer conn.Close()

    duration := time.Since(start)
    return int(duration.Milliseconds()), nil
}

func TestAllLatency(ctx context.Context, servers []model.Server) []model.Server {
    var wg sync.WaitGroup
    results := make([]model.Server, len(servers))
    
    // محدود کردن همزمانی به 50 گوروتین برای جلوگیری از تخلیه File Descriptor
    sem := make(chan struct{}, 50)

    for i, srv := range servers {
        wg.Add(1)
        go func(idx int, server model.Server) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()
            
            latency, err := TestLatency(ctx, server.Address, server.Port)
            if err != nil {
                server.Latency = 9999
            } else {
                server.Latency = latency
            }
            results[idx] = server
        }(i, srv)
    }

    wg.Wait()
    return results
}
