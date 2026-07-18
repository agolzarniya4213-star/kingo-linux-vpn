package core

import (
    "context"
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
    return &SingBoxManager{state: StateDisconnected}
}

func (m *SingBoxManager) Start(ctx context.Context, configPath string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.state == StateConnected || m.state == StateConnecting {
        return nil
    }

    m.state = StateConnecting
    m.cmd = exec.CommandContext(ctx, "sing-box", "run", "-c", configPath)
    m.cmd.Stdout = os.Stdout
    m.cmd.Stderr = os.Stderr

    if err := m.cmd.Start(); err != nil {
        m.state = StateError
        return err
    }

    go func() {
        m.cmd.Wait()
        m.mu.Lock()
        defer m.mu.Unlock()
        if m.state != StateDisconnected {
            m.state = StateDisconnected
        }
    }()

    m.state = StateConnected
    return nil
}

func (m *SingBoxManager) Stop() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.cmd != nil && m.cmd.Process != nil {
        m.cmd.Process.Kill()
    }
    m.cmd = nil
    m.state = StateDisconnected
    return nil
}

func (m *SingBoxManager) GetState() State {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.state
}
