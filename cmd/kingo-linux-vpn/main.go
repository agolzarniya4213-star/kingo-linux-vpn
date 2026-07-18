package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/fetcher"
)

func main() {
    f := fetcher.NewHttpFetcher()
    ctx := context.Background()

    servers, err := f.Fetch(ctx, "https://raw.githubusercontent.com/kingowow/Kingo-vpn/main/merged_config.txt")
    if err != nil {
        log.Fatalf("Failed to fetch servers: %v", err)
    }
    fmt.Printf("Fetched %d servers\n", len(servers))

    ipc := core.NewIpcServer()
    ipc.SetServers(servers)

    if err := ipc.Start("9876"); err != nil {
        log.Fatalf("Failed to start IPC: %v", err)
    }

    fmt.Println("Daemon running. Press Ctrl+C to stop.")

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    <-sigs

    fmt.Println("\nShutting down...")
    ipc.Stop()
    time.Sleep(100 * time.Millisecond)
}
