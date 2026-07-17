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
        fmt.Printf("Failed to start daemon: %v\n", err)
        return
    }
    defer listener.Close()
    defer os.Remove(socketPath)

    fmt.Println("[OK] Mock Daemon is securely running on:", socketPath)

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
            res = Response{Status: "success", Message: "Securely connected to mock server."}
        } else if req.Action == "disconnect" {
            res = Response{Status: "success", Message: "Successfully disconnected."}
        } else {
            res = Response{Status: "error", Message: "Invalid command"}
        }

        json.NewEncoder(conn).Encode(res)
        conn.Close()
        fmt.Println("[EVENT] Handled action:", req.Action)
    }
}
