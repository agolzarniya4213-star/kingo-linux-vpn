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

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/config"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/ipc"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/storage"
)

func getDBPath() string {
    if os.Geteuid() == 0 { return "/var/lib/kingo-vpn/kingo.db" }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".local", "share", "kingo-vpn", "kingo.db")
}

func getKeyPath() string {
    if os.Geteuid() == 0 { return "/var/lib/kingo-vpn/key.bin" }
    home, _ := os.UserHomeDir()
    return filepath.Join(home, ".local", "share", "kingo-vpn", "key.bin")
}

func main() {
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    appCfg := config.LoadConfig()

    dbPath := getDBPath()
    keyPath := getKeyPath()
    
    os.MkdirAll(filepath.Dir(dbPath), 0700)

    crypto, err := storage.NewCryptoLayer(keyPath)
    if err != nil {
        slog.Error("Failed to init crypto layer", "error", err)
        os.Exit(1)
    }

    db, err := storage.NewSQLiteStorage(dbPath, crypto)
    if err != nil {
        slog.Error("Failed to init database", "error", err)
        os.Exit(1)
    }
    defer db.Close()
    os.Chmod(dbPath, 0600)

    // FIX: Seed public test servers if DB is empty (like Hiddify)
    if servers, _ := db.GetServers(); len(servers) == 0 {
        slog.Info("Seeding default public servers...")
        db.SaveServers([]model.Server{
            {ID: "pub1", Name: "Public - Germany (VLESS)", Protocol: "vless", Address: "speedtest.tele2.net", Port: 443, URI: "vless://uuid@speedtest.tele2.net:443?security=tls&type=ws&path=%2F#Germany-Test"},
            {ID: "pub2", Name: "Public - Cloudflare (VLESS)", Protocol: "vless", Address: "1.1.1.1", Port: 443, URI: "vless://uuid@1.1.1.1:443?security=tls&type=ws&path=%2F#Cloudflare-Test"},
        })
    }

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() { <-sigChan; cancel() }()

    coreManager := core.NewSingBoxManager()
    ipcServer := ipc.NewServer(appCfg, coreManager, db)

    go func() {
        if err := ipcServer.Start(ctx); err != nil {
            slog.Error("IPC server failed", "error", err)
            cancel()
        }
    }()

    daemon.SdNotify(false, daemon.SdNotifyReady)
    
    interval, err := daemon.SdWatchdogEnabled(false)
    if err != nil || interval == 0 { interval = 15 * time.Second }
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
