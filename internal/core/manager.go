package core

import (
    "context"
    "encoding/json"
    "fmt"
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
    Start(ctx context.Context, configPath string, clashSecret string) error
    Stop() error
    GetState() State
    GetTraffic() TrafficStats
}

type SingBoxManager struct {
    mu            sync.RWMutex
    state         State
    cmd           *exec.Cmd
    traffic       TrafficStats
    cancel        context.CancelFunc
    configPath    string
    clashSecret   string
    failureCount  int
    lastFailure   time.Time
}

func NewSingBoxManager() *SingBoxManager {
    return &SingBoxManager{state: StateDisconnected}
}

// FIX BUG-036: Proper State Machine transitions
func (m *SingBoxManager) setState(newState State) {
    valid := false
    switch m.state {
    case StateDisconnected:
        if newState == StateConnecting { valid = true }
    case StateConnecting:
        if newState == StateConnected || newState == StateError || newState == StateDisconnected { valid = true }
    case StateConnected:
        if newState == StateDisconnected || newState == StateError { valid = true }
    case StateError:
        if newState == StateConnecting || newState == StateDisconnected { valid = true }
    }
    if valid {
        m.state = newState
    }
}

func (m *SingBoxManager) Start(ctx context.Context, configPath string, clashSecret string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // FIX BUG-037: Circuit Breaker
    if m.failureCount >= 3 && time.Since(m.lastFailure) < 1*time.Minute {
        return fmt.Errorf("circuit breaker tripped: too many failures, please wait")
    }

    if m.state == StateConnected || m.state == StateConnecting {
        return nil
    }

    if _, err := exec.LookPath("sing-box"); err != nil {
        m.setState(StateError)
        return fmt.Errorf("sing-box binary not found in system PATH")
    }

    // Validate config
    checkCmd := exec.Command("sing-box", "check", "-c", configPath)
    if err := checkCmd.Run(); err != nil {
        m.setState(StateError)
        m.failureCount++
        m.lastFailure = time.Now()
        return fmt.Errorf("config validation failed: %w", err)
    }

    m.setState(StateConnecting)
    m.configPath = configPath
    m.clashSecret = clashSecret
    m.cmd = exec.CommandContext(ctx, "sing-box", "run", "-c", configPath)
    
    // FIX BUG-034: Redirect sing-box logs to file instead of stdout to prevent credential leaks
    logFile, err := os.OpenFile("/var/log/kingo-vpn/sing-box.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
    if err != nil {
        m.cmd.Stdout = os.Stdout
        m.cmd.Stderr = os.Stderr
    } else {
        defer logFile.Close()
        m.cmd.Stdout = logFile
        m.cmd.Stderr = logFile
    }

    if err := m.cmd.Start(); err != nil {
        m.setState(StateError)
        m.failureCount++
        m.lastFailure = time.Now()
        return err
    }

    cmd := m.cmd
    trafficCtx, cancel := context.WithCancel(ctx)
    m.cancel = cancel

    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Printf("Recovered from panic: %v\n", r)
            }
        }()
        _ = cmd.Wait()
        m.mu.Lock()
        defer m.mu.Unlock()
        m.setState(StateDisconnected)
        if m.cancel != nil {
            m.cancel()
            m.cancel = nil
        }
        if m.configPath != "" {
            os.Remove(m.configPath)
            m.configPath = ""
        }
    }()

    // FIX BUG-016: Validate actual connection before setting StateConnected
    go m.verifyConnection(trafficCtx)

    go m.monitorTraffic(trafficCtx)
    return nil
}

func (m *SingBoxManager) verifyConnection(ctx context.Context) {
    client := &http.Client{Timeout: 2 * time.Second}
    for i := 0; i < 10; i++ { // Retry for up to 10 seconds
        select {
        case <-ctx.Done():
            return
        case <-time.After(1 * time.Second):
        }
        
        req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:9090/version", nil)
        if m.clashSecret != "" {
            req.Header.Set("Authorization", "Bearer "+m.clashSecret)
        }
        resp, err := client.Do(req)
        if err == nil && resp.StatusCode == 200 {
            resp.Body.Close()
            m.mu.Lock()
            if m.state == StateConnecting {
                m.setState(StateConnected)
                m.failureCount = 0 // Reset on success
            }
            m.mu.Unlock()
            return
        }
    }
    
    // If we reach here, connection failed
    m.mu.Lock()
    m.setState(StateError)
    m.failureCount++
    m.lastFailure = time.Now()
    if m.cancel != nil {
        m.cancel()
    }
    m.mu.Unlock()
}

func (m *SingBoxManager) monitorTraffic(ctx context.Context) {
    time.Sleep(3 * time.Second) // Wait for connection to establish

    m.mu.RLock()
    secret := m.clashSecret
    m.mu.RUnlock()

    req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:9090/traffic", nil)
    if secret != "" {
        req.Header.Set("Authorization", "Bearer "+secret)
    }
    
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
                time.Sleep(1 * time.Second)
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
    m.setState(StateDisconnected)
    m.traffic = TrafficStats{}
    
    if m.configPath != "" {
        os.Remove(m.configPath)
        m.configPath = ""
    }
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
