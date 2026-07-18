package core

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
    "os"
    "os/exec"
    "sync"
    "time"
)

type State string

const (
    StateDisconnected State = "disconnected"
    StateConnecting   State = "connecting"
    StateConnected    State = "connected"
    StateError        State = "error"
)

type TrafficStats struct {
    Upload   int64
    Download int64
}

type Manager interface {
    Start(ctx context.Context, configPath string) error
    Stop() error
    GetState() State
    GetTraffic() TrafficStats
}

type SingBoxManager struct {
    mu      sync.RWMutex
    state   State
    cmd     *exec.Cmd
    traffic TrafficStats
    cancel  context.CancelFunc
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

    // ایجاد یک context جداگانه برای مانیتورینگ ترافیک تا با توقف VPN متوقف شود
    trafficCtx, cancel := context.WithCancel(context.Background())
    m.cancel = cancel

    go func() {
        m.cmd.Wait()
        m.mu.Lock()
        defer m.mu.Unlock()
        if m.state != StateDisconnected {
            m.state = StateDisconnected
        }
        if m.cancel != nil {
            m.cancel()
            m.cancel = nil
        }
    }()

    m.state = StateConnected
    
    // شروع مانیتورینگ ترافیک
    go m.monitorTraffic(trafficCtx)

    return nil
}

func (m *SingBoxManager) monitorTraffic(ctx context.Context) {
    // صبر برای بالا آمدن Clash API
    time.Sleep(2 * time.Second)

    req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:9090/traffic", nil)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    decoder := json.NewDecoder(resp.Body)
    for {
        select {
        case <-ctx.Done():
            return
        default:
            var t struct {
                Up   int64 `json:"up"`
                Down int64 `json:"down"`
            }
            if err := decoder.Decode(&t); err != nil {
                if err == io.EOF {
                    return
                }
                continue
            }
            m.mu.Lock()
            m.traffic = TrafficStats{Upload: t.Up, Download: t.Down}
            m.mu.Unlock()
        }
    }
}

func (m *SingBoxManager) Stop() error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.cancel != nil {
        m.cancel()
        m.cancel = nil
    }

    if m.cmd != nil && m.cmd.Process != nil {
        m.cmd.Process.Kill()
    }
    m.cmd = nil
    m.state = StateDisconnected
    m.traffic = TrafficStats{}
    return nil
}

func (m *SingBoxManager) GetState() State {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.state
}

func (m *SingBoxManager) GetTraffic() TrafficStats {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.traffic
}
