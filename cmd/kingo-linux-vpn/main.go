package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/daemon"
)

func main() {
    configDir, err := os.UserConfigDir()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error getting config dir: %v\n", err)
        os.Exit(1)
    }

    ipcServer := daemon.NewIpcServer(configDir)
    if err := ipcServer.Start(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to start daemon: %v\n", err)
        os.Exit(1)
    }
    defer ipcServer.Stop()

    fmt.Println("Kingo Daemon started.")
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
    fmt.Println("Shutting down daemon...")
}
