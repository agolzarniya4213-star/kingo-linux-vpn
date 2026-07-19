package core

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
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

func (m *SingBoxManager) setState(newState State) {
    m.state = newState
}

func findSingBox() (string, error) {
    if path, err := exec.LookPath("sing-box"); err == nil {
        return path, nil
    }
    
    candidates := []string{"./sing-box", "./build/sing-box", "../sing-box", "../build/sing-box"}
    if exe, err := os.Executable(); err == nil {
        exeDir := filepath.Dir(exe)
        candidates = append(candidates, filepath.Join(exeDir, "sing-box"), filepath.Join(exeDir, "build", "sing-box"))
    }
    
    for _, c := range candidates {
        if abs, err := filepath.Abs(c); err == nil {
            if _, err := os.Stat(abs); err == nil {
                return abs, nil
            }
        }
    }
    return "", fmt.Errorf("sing-box binary not found")
}

func (m *SingBoxManager) Start(ctx context.Context, configPath string, clashSecret string) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    if m.failureCount >= 3 && time.Since(m.lastFailure) < 1*time.Minute {
        return fmt.Errorf("circuit breaker tripped: too many failures")
    }

    if m.state == StateConnected || m.state == StateConnecting {
        return nil
    }

    singBoxPath, err := findSingBox()
    if err != nil {
        m.setState(StateError)
        return err
    }

    checkCmd := exec.Command(singBoxPath, "check", "-c", configPath)
    if err := checkCmd.Run(); err != nil {
        m.setState(StateError)
        m.failureCount++
        m.lastFailure = time.Now()
        return fmt.Errorf("invalid config: %w", err)
    }

    m.setState(StateConnecting)
    m.configPath = configPath
    m.clashSecret = clashSecret
    m.cmd = exec.CommandContext(ctx, singBoxPath, "run", "-c", configPath)
    
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
        defer func() { if r := recover(); r != nil { fmt.Printf("Recovered: %v\n", r) } }()
        _ = cmd.Wait()
        m.mu.Lock()
        defer m.mu.Unlock()
        m.setState(StateDisconnected)
        if m.cancel != nil { m.cancel(); m.cancel = nil }
        if m.configPath != "" { os.Remove(m.configPath); m.configPath = "" }
    }()

    go m.verifyConnection(trafficCtx)
    go m.monitorTraffic(trafficCtx)
    return nil
}

func (m *SingBoxManager) verifyConnection(ctx context.Context) {
    client := &http.Client{Timeout: 2 * time.Second}
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done(): return
        case <-time.After(1 * time.Second):
        }
        req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:9090/version", nil)
        if m.clashSecret != "" { req.Header.Set("Authorization", "Bearer "+m.clashSecret) }
        resp, err := client.Do(req)
        if err == nil && resp.StatusCode == 200 {
            resp.Body.Close()
            m.mu.Lock()
            if m.state == StateConnecting { m.setState(StateConnected); m.failureCount = 0 }
            m.mu.Unlock()
            return
        }
    }
    m.mu.Lock()
    m.setState(StateError)
    m.failureCount++
    m.lastFailure = time.Now()
    if m.cancel != nil { m.cancel() }
    m.mu.Unlock()
}

func (m *SingBoxManager) monitorTraffic(ctx context.Context) {
    time.Sleep(3 * time.Second)
    m.mu.RLock()
    secret := m.clashSecret
    m.mu.RUnlock()

    req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:9090/traffic", nil)
    if secret != "" { req.Header.Set("Authorization", "Bearer "+secret) }
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil { return }
    defer resp.Body.Close()

    decoder := json.NewDecoder(resp.Body)
    for {
        select {
        case <-ctx.Done(): return
        default:
            var t struct { Up int64 `json:"up"`; Down int64 `json:"down"` }
            if err := decoder.Decode(&t); err != nil {
                if err == io.EOF { return }
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
    if m.cancel != nil { m.cancel(); m.cancel = nil }
    if m.cmd != nil && m.cmd.Process != nil { m.cmd.Process.Kill() }
    m.setState(StateDisconnected)
    m.traffic = TrafficStats{}
    if m.configPath != "" { os.Remove(m.configPath); m.configPath = "" }
    return nil
}

func (m *SingBoxManager) GetState() State { m.mu.RLock(); defer m.mu.RUnlock(); return m.state }
func (m *SingBoxManager) GetTraffic() TrafficStats { m.mu.RLock(); defer m.mu.RUnlock(); return m.traffic }
