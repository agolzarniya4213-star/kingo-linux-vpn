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
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    coreManager := core.NewSingBoxManager()
    ipcServer := ipc.NewServer("/tmp/kingo-vpn.sock", coreManager)

    if err := ipcServer.Start(ctx); err != nil {
        slog.Error("IPC server failed", "error", err)
        os.Exit(1)
    }
}
