package ipc

import (
    "context"
    "encoding/json"
    "io"
    "net"
    "os"
    "path/filepath"
    "sort"
    "time"

    "golang.org/x/sys/unix"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/config"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/fetcher"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/network"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/storage"
)

const ProtocolVersion = "2.0"

type Request struct {
    RequestID  string `json:"request_id"`
    Version    string `json:"version"`
    Action     string `json:"action"`
    ConfigPath string `json:"config_path,omitempty"`
    SubURL     string `json:"sub_url,omitempty"`
    ServerURI  string `json:"server_uri,omitempty"`
}

type Response struct {
    RequestID string `json:"request_id"`
    Success   bool           `json:"success"`
    Message   string         `json:"message"`
    State     string         `json:"state"`
    Servers   []model.Server `json:"servers,omitempty"`
    Upload    int64          `json:"upload,omitempty"`
    Download  int64          `json:"download,omitempty"`
}

type Server struct {
    appCfg  *config.AppConfig
    manager core.Manager
    db      *storage.SQLiteStorage
}

func NewServer(appCfg *config.AppConfig, manager core.Manager, db *storage.SQLiteStorage) *Server {
    return &Server{appCfg: appCfg, manager: manager, db: db}
}

func (s *Server) Start(ctx context.Context) error {
    socketDir := "/run/kingo-vpn"
    if os.Geteuid() != 0 { socketDir = "/tmp/kingo-vpn" }
    os.MkdirAll(socketDir, 0755)
    s.appCfg.IPC.SocketPath = filepath.Join(socketDir, "kingo-vpn.sock")

    if _, err := os.Stat(s.appCfg.IPC.SocketPath); err == nil { os.Remove(s.appCfg.IPC.SocketPath) }
    listener, err := net.Listen("unix", s.appCfg.IPC.SocketPath)
    if err != nil { return err }
    os.Chmod(s.appCfg.IPC.SocketPath, 0660)
    defer listener.Close()

    go func() { <-ctx.Done(); listener.Close(); os.Remove(s.appCfg.IPC.SocketPath) }()
    sem := make(chan struct{}, s.appCfg.IPC.MaxConns)

    for {
        conn, err := listener.Accept()
        if err != nil { select { case <-ctx.Done(): return nil; default: continue } }
        select {
        case sem <- struct{}{}:
            go func() { defer func() { <-sem }(); s.handleConnection(conn) }()
        default:
            conn.Write([]byte(`{"success":false,"message":"server busy"}`)); conn.Close()
        }
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    uc, ok := conn.(*net.UnixConn)
    if !ok { return }
    rawConn, err := uc.SyscallConn()
    if err != nil { return }
    var uid int = -1
    rawConn.Control(func(fd uintptr) { ucred, _ := unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED); uid = int(ucred.Uid) })
    if uid != 0 && uid != os.Getuid() { conn.Write([]byte(`{"success":false,"message":"unauthorized"}`)); return }

    decoder := json.NewDecoder(io.LimitReader(conn, 1<<20))
    encoder := json.NewEncoder(conn)

    for {
        conn.SetReadDeadline(time.Now().Add(10 * time.Second))
        var req Request
        if err := decoder.Decode(&req); err != nil { break }
        conn.SetReadDeadline(time.Time{})

        var resp Response
        resp.RequestID = req.RequestID
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        
        switch req.Action {
        case "connect_server":
            configPath, clashSecret, err := config.GenerateSingBoxConfig(req.ServerURI, s.appCfg)
            if err != nil {
                resp = Response{RequestID: req.RequestID, Success: false, Message: "Config Error: " + err.Error()}
                break
            }
            err = s.manager.Start(ctx, configPath, clashSecret)
            resp = Response{RequestID: req.RequestID, Success: err == nil, State: string(s.manager.GetState())}
            if err != nil { resp.Message = "VPN Start Failed: " + err.Error() } // Return actual error to UI

        case "auto_connect":
            servers, err := s.db.GetServers()
            if err != nil || len(servers) == 0 {
                resp = Response{RequestID: req.RequestID, Success: false, Message: "No servers available"}
                break
            }
            testedServers := network.TestAllLatency(ctx, servers)
            _ = s.db.SaveServers(testedServers)
            sort.Slice(testedServers, func(i, j int) bool { return testedServers[i].Latency < testedServers[j].Latency })
            
            var bestServer model.Server
            found := false
            for _, srv := range testedServers { if srv.Latency < 9999 { bestServer = srv; found = true; break } }
            if !found { resp = Response{RequestID: req.RequestID, Success: false, Message: "All unreachable", Servers: testedServers}; break }
            
            configPath, clashSecret, err := config.GenerateSingBoxConfig(bestServer.URI, s.appCfg)
            if err != nil { resp = Response{RequestID: req.RequestID, Success: false, Message: "Config gen failed", Servers: testedServers}; break }
            err = s.manager.Start(ctx, configPath, clashSecret)
            resp = Response{RequestID: req.RequestID, Success: err == nil, State: string(s.manager.GetState()), Servers: testedServers}
            if err != nil { resp.Message = "VPN Start Failed: " + err.Error() }

        case "disconnect":
            s.manager.Stop()
            resp = Response{RequestID: req.RequestID, Success: true, State: string(s.manager.GetState())}
        case "status":
            resp = Response{RequestID: req.RequestID, Success: true, State: string(s.manager.GetState())}
        case "get_servers":
            servers, err := s.db.GetServers()
            resp = Response{RequestID: req.RequestID, Success: err == nil, Servers: servers}
        case "add_subscription":
            servers, err := fetcher.FetchSubscription(req.SubURL)
            if err != nil { resp = Response{RequestID: req.RequestID, Success: false, Message: "Fetch failed"} } else {
                err = s.db.SaveServers(servers)
                resp = Response{RequestID: req.RequestID, Success: err == nil, Servers: servers}
            }
        case "clear_servers":
            // FIX: Added clear servers action
            err := s.db.SaveServers([]model.Server{})
            resp = Response{RequestID: req.RequestID, Success: err == nil, Servers: []model.Server{}}
        case "test_latency":
            servers, err := s.db.GetServers()
            if err != nil { resp = Response{RequestID: req.RequestID, Success: false}; break }
            testedServers := network.TestAllLatency(ctx, servers)
            _ = s.db.SaveServers(testedServers)
            resp = Response{RequestID: req.RequestID, Success: true, Servers: testedServers}
        case "get_traffic":
            traffic := s.manager.GetTraffic()
            resp = Response{RequestID: req.RequestID, Success: true, Upload: traffic.Upload, Download: traffic.Download}
        default:
            resp = Response{RequestID: req.RequestID, Success: false, Message: "unknown action"}
        }
        cancel()
        if err := encoder.Encode(resp); err != nil { break }
    }
}
