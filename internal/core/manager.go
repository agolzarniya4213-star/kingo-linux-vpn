package core

import (
    "context"
    "log/slog"
    "os"
    "os/exec"
    "sync"
)

type State string

const (
    StateDisconnected State = "disconnected"
    StateConnecting   State = "connecting"
    StateConnected    State = "connected"
    StateError        State = "error"
)

type Manager interface {
    Start(ctx context.Context, configPath string) error
    Stop() error
    GetState() State
}

type SingBoxManager struct {
    mu    sync.RWMutex
    state State
    cmd   *exec.Cmd
}

func NewSingBoxManager() *SingBoxManager {
    return &SingBoxManager{
        state: StateDisconnected,
    }
}

func (m *SingBoxManager) Start(ctx context.Context, configPath string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.state == StateConnected || m.state == StateConnecting {
        slog.Warn("VPN is already running or connecting")
        return nil
    }

    m.state = StateConnecting
    slog.Info("Starting Sing-box core...", "config", configPath)

    m.cmd = exec.CommandContext(ctx, "sing-box", "run", "-c", configPath)
    m.cmd.Stdout = os.Stdout
    m.cmd.Stderr = os.Stderr
    
    if err := m.cmd.Start(); err != nil {
        m.state = StateError
        slog.Error("Failed to start Sing-box", "error", err)
        return err
    }

    // گوروتین برای مانیتورینگ پایان کار پردازش
    go func() {
        err := m.cmd.Wait()
        m.mu.Lock()
        defer m.mu.Unlock()
        
        if m.state != StateDisconnected {
            if err != nil {
                slog.Error("Sing-box process exited with error", "error", err)
                m.state = StateError
            } else {
                slog.Info("Sing-box process exited successfully")
                m.state = StateDisconnected
            }
        }
    }()

    m.state = StateConnected
    slog.Info("Sing-box started successfully")
    return nil
}

func (m *SingBoxManager) Stop() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.cmd != nil && m.cmd.Process != nil {
        slog.Info("Stopping Sing-box core...")
        err := m.cmd.Process.Kill()
        if err != nil {
            slog.Error("Failed to kill Sing-box process", "error", err)
            m.state = StateError
            return err
        }
    }

    m.cmd = nil
    m.state = StateDisconnected
    slog.Info("Sing-box stopped successfully")
    return nil
}

func (m *SingBoxManager) GetState() State {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.state
}
