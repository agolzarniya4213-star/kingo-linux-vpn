package ipc

import (
    "context"
    "encoding/json"
    "net"
    "os"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/model"
    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/storage"
)

type Request struct {
    Action     string `json:"action"`
    ConfigPath string `json:"config_path,omitempty"`
}

type Response struct {
    Success bool           `json:"success"`
    Message string         `json:"message"`
    State   string         `json:"state"`
    Servers []model.Server `json:"servers,omitempty"`
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
        case "connect":
            err := s.manager.Start(context.Background(), req.ConfigPath)
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
        default:
            resp = Response{Success: false, Message: "unknown action"}
        }
        encoder.Encode(resp)
    }
}
