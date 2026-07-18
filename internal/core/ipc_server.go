package core

import (
    "encoding/json"
    "fmt"
    "log"
    "net"
    "strings"
    "sync"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
)

type IpcServer struct {
    listener net.Listener
    servers  []model.Server
    mu       sync.RWMutex
}

func NewIpcServer() *IpcServer {
    return &IpcServer{}
}

func (s *IpcServer) Start(port string) error {
    var err error
    s.listener, err = net.Listen("tcp", ":"+port)
    if err != nil {
        return fmt.Errorf("listen failed: %w", err)
    }
    log.Printf("IPC server listening on %s", s.listener.Addr().String())
    go s.acceptLoop()
    return nil
}

func (s *IpcServer) SetServers(servers []model.Server) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.servers = servers
}

func (s *IpcServer) Stop() {
    if s.listener != nil {
        s.listener.Close()
    }
}

func (s *IpcServer) acceptLoop() {
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
    buf := make([]byte, 4096)
    n, err := conn.Read(buf)
    if err != nil {
        return
    }
    cmd := strings.TrimSpace(string(buf[:n]))

    var response []byte
    switch cmd {
    case "refresh":
        s.mu.RLock()
        servers := s.servers
        s.mu.RUnlock()
        resp := map[string]interface{}{"servers": servers}
        response, _ = json.Marshal(resp)
    case "status":
        resp := map[string]string{"status": "Disconnected"}
        response, _ = json.Marshal(resp)
    default:
        resp := map[string]string{"error": "unknown command"}
        response, _ = json.Marshal(resp)
    }
    conn.Write(response)
    conn.Write([]byte("\n"))
}
