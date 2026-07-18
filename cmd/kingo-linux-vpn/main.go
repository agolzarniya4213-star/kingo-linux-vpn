package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/fetcher"
)

func main() {
    f := fetcher.NewHttpFetcher()
    ctx := context.Background()

    servers, err := f.Fetch(ctx, "https://raw.githubusercontent.com/kingowow/Kingo-vpn/main/merged_config.txt")
    if err != nil {
        log.Fatalf("Failed to fetch servers: %v", err)
    }

    fmt.Printf("Successfully parsed %d servers:\n", len(servers))
    for i, s := range servers {
        fmt.Printf("%d. [%s] %s:%d - %s\n", i+1, s.Protocol, s.Address, s.Port, s.Name)
    }
    os.Exit(0)
}
