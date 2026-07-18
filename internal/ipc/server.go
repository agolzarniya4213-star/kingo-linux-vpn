package ipc

import (
    "context"
    "encoding/json"
    "net"
    "os"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/config"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/fetcher"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/network"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/storage"
)

type Request struct {
    Action     string `json:"action"`
    ConfigPath string `json:"config_path,omitempty"`
    SubURL     string `json:"sub_url,omitempty"`
    ServerURI  string `json:"server_uri,omitempty"`
}

type Response struct {
    Success  bool           `json:"success"`
    Message  string         `json:"message"`
    State    string         `json:"state"`
    Servers  []model.Server `json:"servers,omitempty"`
    Upload   int64          `json:"upload,omitempty"`
    Download int64          `json:"download,omitempty"`
}

type Server struct {
    socketPath string
    manager    core.Manager
    db         *storage.SQLiteStorage
}

func NewServer(socketPath string, manager core.Manager, db *storage.SQLiteStorage) *Server {
    return &Server{socketPath: socketPath, manager: manager, db: db}
}

func (s *Server) Start(ctx context.Context) error {
    if _, err := os.Stat(s.socketPath); err == nil {
        os.Remove(s.socketPath)
    }

    listener, err := net.Listen("unix", s.socketPath)
    if err != nil {
        return err
    }
    defer listener.Close()

    go func() {
        <-ctx.Done()
        listener.Close()
        os.Remove(s.socketPath)
    }()

    for {
        conn, err := listener.Accept()
        if err != nil {
            select {
            case <-ctx.Done():
                return nil
            default:
                continue
            }
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    decoder := json.NewDecoder(conn)
    encoder := json.NewEncoder(conn)

    for {
        var req Request
        if err := decoder.Decode(&req); err != nil {
            break
        }

        var resp Response
        switch req.Action {
        case "connect_server":
            configPath, err := config.GenerateSingBoxConfig(req.ServerURI)
            if err != nil {
                resp = Response{Success: false, Message: "Config gen failed: " + err.Error()}
                break
            }
            err = s.manager.Start(context.Background(), configPath)
            resp = Response{Success: err == nil, State: string(s.manager.GetState())}
            if err != nil {
                resp.Message = err.Error()
            }
        case "disconnect":
            s.manager.Stop()
            resp = Response{Success: true, State: string(s.manager.GetState())}
        case "status":
            resp = Response{Success: true, State: string(s.manager.GetState())}
        case "get_servers":
            servers, err := s.db.GetServers()
            resp = Response{Success: err == nil, Servers: servers}
            if err != nil {
                resp.Message = err.Error()
            }
        case "add_subscription":
            servers, err := fetcher.FetchSubscription(req.SubURL)
            if err != nil {
                resp = Response{Success: false, Message: err.Error()}
            } else {
                err = s.db.SaveServers(servers)
                resp = Response{Success: err == nil, Servers: servers}
                if err != nil {
                    resp.Message = err.Error()
                }
            }
        case "test_latency":
            servers, err := s.db.GetServers()
            if err != nil {
                resp = Response{Success: false, Message: err.Error()}
                break
            }
            testedServers := network.TestAllLatency(context.Background(), servers)
            _ = s.db.SaveServers(testedServers)
            resp = Response{Success: true, Servers: testedServers}
        case "get_traffic":
            traffic := s.manager.GetTraffic()
            resp = Response{Success: true, Upload: traffic.Upload, Download: traffic.Download}
        default:
            resp = Response{Success: false, Message: "unknown action"}
        }
        encoder.Encode(resp)
    }
}
