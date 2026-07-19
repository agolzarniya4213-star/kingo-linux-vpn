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
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/storage"
)

func getDBPath() string {
    // اگر سرویس با دسترسی Root اجرا شود (systemd)
    if os.Geteuid() == 0 {
        return "/var/lib/kingo-vpn/kingo.db"
    }
    // اگر توسط کاربر عادی اجرا شود (run.sh)
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

    // اطمینان از وجود پوشه دیتابیس
    dbPath := getDBPath()
    os.MkdirAll(filepath.Dir(dbPath), 0755)

    db, err := storage.NewSQLiteStorage(dbPath)
    if err != nil {
        slog.Error("Failed to init database", "error", err)
        os.Exit(1)
    }
    defer db.Close()

    if s, _ := db.GetServers(); len(s) == 0 {
        db.SaveServers([]model.Server{
            {ID: "srv1", Name: "Germany - Frankfurt", Protocol: "vless", Address: "10.0.0.1", Port: 443},
            {ID: "srv2", Name: "Netherlands - Amsterdam", Protocol: "vmess", Address: "10.0.0.2", Port: 8080},
        })
        slog.Info("Seeded dummy servers to database")
    }

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()

    coreManager := core.NewSingBoxManager()
    ipcServer := ipc.NewServer("/tmp/kingo-vpn.sock", coreManager, db)

    if err := ipcServer.Start(ctx); err != nil {
        slog.Error("IPC server failed", "error", err)
        os.Exit(1)
    }
}
