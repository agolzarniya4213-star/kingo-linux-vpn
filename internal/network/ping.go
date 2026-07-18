package network

import (
    "context"
    "net"
    "strconv"
    "sync"
    "time"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

// TestLatency برای یک سرور تأخیر را تست می‌کند (با پشتیبانی صحیح از IPv6)
func TestLatency(ctx context.Context, host string, port int) (int, error) {
    // استفاده از JoinHostPort برای پشتیبانی صحیح از IPv6
    address := net.JoinHostPort(host, strconv.Itoa(port))
    start := time.Now()

    // ایجاد یک تایم‌اوت ۳ ثانیه‌ای از طریق Context
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    d := net.Dialer{}
    conn, err := d.DialContext(ctx, "tcp", address)
    if err != nil {
        return 0, err
    }
    defer conn.Close()

    // زمان صرف شده برای TCP Handshake
    duration := time.Since(start)
    return int(duration.Milliseconds()), nil
}

// TestAllLatency تأخیر تمام سرورها را به صورت همزمان تست می‌کند
func TestAllLatency(ctx context.Context, servers []model.Server) []model.Server {
    var wg sync.WaitGroup
    results := make([]model.Server, len(servers))

    for i, srv := range servers {
        wg.Add(1)
        go func(idx int, server model.Server) {
            defer wg.Done()
            
            latency, err := TestLatency(ctx, server.Address, server.Port)
            if err != nil {
                server.Latency = 9999 // نشان‌دهنده خطا یا تایم‌اوت
            } else {
                server.Latency = latency
            }
            results[idx] = server
        }(i, srv)
    }

    wg.Wait()
    return results
}
