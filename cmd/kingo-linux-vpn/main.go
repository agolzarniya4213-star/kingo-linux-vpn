package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/ipc"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    slog.Info("Starting Kingo Linux VPN Daemon...")

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        sig := <-sigChan
        slog.Info("Received signal, shutting down...", "signal", sig)
        cancel()
    }()

    coreManager := core.NewSingBoxManager()
    ipcServer := ipc.NewServer("/tmp/kingo-vpn.sock", coreManager)

    err := ipcServer.Start(ctx)
    if err != nil {
        slog.Error("Failed to start IPC server", "error", err)
        os.Exit(1)
    }

    slog.Info("Kingo Daemon stopped gracefully.")
}
