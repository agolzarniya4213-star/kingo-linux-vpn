package main

import (
    "context"
    "log/slog"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"

    "github.com/coreos/go-systemd/v22/daemon"

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
    dbDir := filepath.Dir(dbPath)
    
    // FIX BUG-006 & 048: Secure directory and file permissions
    if err := os.MkdirAll(dbDir, 0700); err != nil {
        slog.Error("Failed to create secure database directory", "error", err)
        os.Exit(1)
    }

    db, err := storage.NewSQLiteStorage(dbPath)
    if err != nil {
        slog.Error("Failed to init database", "error", err)
        os.Exit(1)
    }
    defer db.Close()
    os.Chmod(dbPath, 0600) // Restrict DB access

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    coreManager := core.NewSingBoxManager()
    ipcServer := ipc.NewServer("/run/kingo-vpn/kingo-vpn.sock", coreManager, db)

    // Start IPC server in background
    go func() {
        if err := ipcServer.Start(ctx); err != nil {
            slog.Error("IPC server failed", "error", err)
            cancel()
        }
    }()

    // FIX BUG-025: Systemd Type=notify and Watchdog support
    daemon.SdNotify(false, daemon.SdNotifyReady)
    
    // Watchdog ticker
    interval, err := daemon.SdWatchdogEnabled(false)
    if err != nil || interval == 0 {
        interval = 15 * time.Second // Fallback if not running under systemd
    }
    ticker := time.NewTicker(interval / 2)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            daemon.SdNotify(false, daemon.SdNotifyStopping)
            return
        case <-ticker.C:
            daemon.SdNotify(false, daemon.SdNotifyWatchdog)
        }
    }
}
