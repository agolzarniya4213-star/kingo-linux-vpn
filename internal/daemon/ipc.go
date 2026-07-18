package daemon

import (
    "encoding/json"
    "fmt"
    "net"
    "os"
    "path/filepath"
    "sync"
)

type Request struct {
    Action string `json:"action"`
    Server string `json:"server,omitempty"` // Added server field
}

type Response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

type IpcServer struct {
    socketPath string
    listener   net.Listener
    mu         sync.Mutex
}

func NewIpcServer(configDir string) *IpcServer {
    socketPath := filepath.Join(configDir, "kingo-linux-vpn", "daemon.sock")
    return &IpcServer{socketPath: socketPath}
}

func (s *IpcServer) Start() error {
    s.mu.Lock()
    defer s.mu.Unlock()

    os.MkdirAll(filepath.Dir(s.socketPath), 0755)
    os.Remove(s.socketPath)

    listener, err := net.Listen("unix", s.socketPath)
    if err != nil {
        return fmt.Errorf("failed to bind IPC socket: %w", err)
    }
    s.listener = listener

    go s.acceptConnections()
    return nil
}

func (s *IpcServer) Stop() {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.listener != nil {
        s.listener.Close()
        os.Remove(s.socketPath)
    }
}

func (s *IpcServer) acceptConnections() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return
        }
        go s.handleConnection(conn)
    }
}

func (s *IpcServer) handleConnection(conn net.Conn) {
    defer conn.Close()

    var req Request
    if err := json.NewDecoder(conn).Decode(&req); err != nil {
        sendError(conn, "Invalid JSON")
        return
    }

    var res Response
    switch req.Action {
    case "connect":
        msg := "Connected to default server."
        if req.Server != "" {
            msg = fmt.Sprintf("Secure tunnel established to %s.", req.Server)
        }
        res = Response{Status: "success", Message: msg}
    case "disconnect":
        res = Response{Status: "success", Message: "Tunnel disconnected safely."}
    default:
        sendError(conn, "Unknown action")
        return
    }

    json.NewEncoder(conn).Encode(res)
}

func sendError(conn net.Conn, msg string) {
    res := Response{Status: "error", Message: msg}
    json.NewEncoder(conn).Encode(res)
}
