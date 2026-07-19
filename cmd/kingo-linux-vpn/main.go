package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/ipc"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/storage"
)

func getDBPath() string {
    if os.Geteuid() == 0 {
        return "/var/lib/kingo-vpn/kingo.db"
    }
    home, err := os.UserHomeDir()
    if err != nil {
        return "kingo.db"
    }
    return filepath.Join(home, ".local", "share", "kingo-vpn", "kingo.db")
}

func main() {
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    dbPath := getDBPath()
    // FIX BUG-048: Check MkdirAll error
    if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
        slog.Error("Failed to create database directory", "error", err)
        os.Exit(1)
    }

    db, err := storage.NewSQLiteStorage(dbPath)
    if err != nil {
        slog.Error("Failed to init database", "error", err)
        os.Exit(1)
    }
    defer db.Close()

    // FIX BUG-007: Removed dummy seed servers for production readiness

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    coreManager := core.NewSingBoxManager()
    // Socket path is now hardcoded inside IPC server to /run/kingo-vpn for security
    ipcServer := ipc.NewServer("/run/kingo-vpn/kingo-vpn.sock", coreManager, db)

    if err := ipcServer.Start(ctx); err != nil {
        slog.Error("IPC server failed", "error", err)
        os.Exit(1)
    }
}
