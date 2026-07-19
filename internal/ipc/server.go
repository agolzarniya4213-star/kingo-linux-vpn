package ipc

import (
    "context"
    "encoding/json"
    "io"
    "net"
    "os"
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
    socketDir := "/run/kingo-vpn"
    os.MkdirAll(socketDir, 0755)
    s.socketPath = socketDir + "/kingo-vpn.sock"

    if _, err := os.Stat(s.socketPath); err == nil {
        os.Remove(s.socketPath)
    }

    listener, err := net.Listen("unix", s.socketPath)
    if err != nil {
        return err
    }
    os.Chmod(s.socketPath, 0600)

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

    uc, ok := conn.(*net.UnixConn)
    if !ok {
        return
    }
    rawConn, err := uc.SyscallConn()
    if err != nil {
        return
    }
    var uid int = -1
    rawConn.Control(func(fd uintptr) {
        ucred, err := unix.GetsockoptUcred(int(fd), unix.SOL_SOCKET, unix.SO_PEERCRED)
        if err == nil {
            uid = int(ucred.Uid)
        }
    })
    if uid != 0 {
        conn.Write([]byte(`{"success":false,"message":"unauthorized"}`))
        return
    }

    decoder := json.NewDecoder(io.LimitReader(conn, 1<<20))
    encoder := json.NewEncoder(conn)

    for {
        conn.SetReadDeadline(time.Now().Add(10 * time.Second))
        var req Request
        if err := decoder.Decode(&req); err != nil {
            break
        }
        conn.SetReadDeadline(time.Time{})

        var resp Response
        switch req.Action {
        case "connect_server":
            configPath, clashSecret, err := config.GenerateSingBoxConfig(req.ServerURI)
            if err != nil {
                resp = Response{Success: false, Message: "Config gen failed: " + err.Error()}
                break
            }
            err = s.manager.Start(context.Background(), configPath, clashSecret)
            resp = Response{Success: err == nil, State: string(s.manager.GetState())}
            if err != nil {
                resp.Message = err.Error()
            }

        case "auto_connect":
            servers, err := s.db.GetServers()
            if err != nil || len(servers) == 0 {
                resp = Response{Success: false, Message: "No servers available"}
                break
            }
            
            testedServers := network.TestAllLatency(context.Background(), servers)
            _ = s.db.SaveServers(testedServers)
            
            sort.Slice(testedServers, func(i, j int) bool {
                return testedServers[i].Latency < testedServers[j].Latency
            })
            
            var bestServer model.Server
            found := false
            for _, srv := range testedServers {
                if srv.Latency < 9999 {
                    bestServer = srv
                    found = true
                    break
                }
            }
            
            if !found {
                resp = Response{Success: false, Message: "All servers are unreachable", Servers: testedServers}
                break
            }
            
            configPath, clashSecret, err := config.GenerateSingBoxConfig(bestServer.URI)
            if err != nil {
                resp = Response{Success: false, Message: "Config gen failed: " + err.Error(), Servers: testedServers}
                break
            }
            err = s.manager.Start(context.Background(), configPath, clashSecret)
            resp = Response{Success: err == nil, State: string(s.manager.GetState()), Servers: testedServers}
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
