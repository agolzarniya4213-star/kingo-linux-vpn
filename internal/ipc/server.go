package ipc

import (
    "context"
    "encoding/json"
    "log/slog"
    "net"
    "os"

    "github.com/agolzarniya4213-star/kingo-linux-vpn/internal/core"
)

type Request struct {
    Action     string `json:"action"`
    ConfigPath string `json:"config_path,omitempty"`
}

type Response struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    State   string `json:"state"`
}

type Server struct {
    socketPath string
    manager    core.Manager
}

func NewServer(socketPath string, manager core.Manager) *Server {
    return &Server{
        socketPath: socketPath,
        manager:    manager,
    }
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

    slog.Info("IPC Server listening", "socket", s.socketPath)

    go func() {
        <-ctx.Done()
        slog.Info("Closing IPC listener...")
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
                slog.Error("Failed to accept connection", "error", err)
                continue
            }
        }
        go s.handleConnection(ctx, conn)
    }
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
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
            err := s.manager.Start(ctx, req.ConfigPath)
            resp = Response{
                Success: err == nil,
                Message: "Connected",
                State:   string(s.manager.GetState()),
            }
            if err != nil {
                resp.Message = err.Error()
            }
        case "disconnect":
            err := s.manager.Stop()
            resp = Response{
                Success: err == nil,
                Message: "Disconnected",
                State:   string(s.manager.GetState()),
            }
        case "status":
            resp = Response{
                Success: true,
                State:   string(s.manager.GetState()),
            }
        default:
            resp = Response{Success: false, Message: "Unknown action"}
        }

        if err := encoder.Encode(resp); err != nil {
            break
        }
    }
}
