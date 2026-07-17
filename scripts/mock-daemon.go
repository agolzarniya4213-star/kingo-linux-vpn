package main

import (
    "encoding/json"
    "fmt"
    "net"
    "os"
    "path/filepath"
)

type Request struct {
    Action string `json:"action"`
}

type Response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

func main() {
    configDir, _ := os.UserConfigDir()
    socketPath := filepath.Join(configDir, "kingo-linux-vpn", "daemon.sock")
    
    os.MkdirAll(filepath.Dir(socketPath), 0755)
    os.Remove(socketPath)

    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    defer os.Remove(socketPath)

    fmt.Println("Mock Daemon is running on:", socketPath)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        
        var req Request
        if err := json.NewDecoder(conn).Decode(&req); err != nil {
            conn.Close()
            continue
        }

        var res Response
        if req.Action == "connect" {
            res = Response{Status: "success", Message: "Mock VPN Connected!"}
        } else if req.Action == "disconnect" {
            res = Response{Status: "success", Message: "Mock VPN Disconnected!"}
        } else {
            res = Response{Status: "error", Message: "Unknown action"}
        }

        json.NewEncoder(conn).Encode(res)
        conn.Close()
        fmt.Println("Handled action:", req.Action)
    }
}
