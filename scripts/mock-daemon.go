package main

import (
    "encoding/json"
    "fmt"
    "net"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
)

// Strict JSON structures to prevent injection/unexpected behavior
type Request struct {
    Action string `json:"action"`
}

type Response struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

func main() {
    configDir, err := os.UserConfigDir()
    if err != nil {
        panic(err)
    }
    
    socketPath := filepath.Join(configDir, "kingo-linux-vpn", "daemon.sock")
    
    // Ensure directory exists
    if err := os.MkdirAll(filepath.Dir(socketPath), 0755); err != nil {
        panic(err)
    }
    
    // Clean up old socket if it exists
    os.Remove(socketPath)

    listener, err := net.Listen("unix", socketPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "FATAL: Cannot bind socket: %v\n", err)
        os.Exit(1)
    }

    // Graceful shutdown handling (NASA level reliability)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        fmt.Println("\n[Mock Daemon] Shutting down gracefully...")
        listener.Close()
        os.Remove(socketPath)
        os.Exit(0)
    }()

    defer listener.Close()
    defer os.Remove(socketPath)

    fmt.Println("[Mock Daemon] Successfully started on:", socketPath)
    fmt.Println("[Mock Daemon] Waiting for UI connections... (Press Ctrl+C to stop)")

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    var req Request
    decoder := json.NewDecoder(conn)
    if err := decoder.Decode(&req); err != nil {
        sendError(conn, "Invalid JSON format")
        return
    }

    var res Response

    switch req.Action {
    case "connect":
        res = Response{Status: "success", Message: "Secure tunnel established successfully."}
    case "disconnect":
        res = Response{Status: "success", Message: "Tunnel disconnected safely."}
    default:
        sendError(conn, "Unauthorized action: "+req.Action)
        return
    }

    json.NewEncoder(conn).Encode(res)
    fmt.Printf("[Mock Daemon] Handled: %s\n", req.Action)
}

func sendError(conn net.Conn, msg string) {
    res := Response{Status: "error", Message: msg}
    json.NewEncoder(conn).Encode(res)
}
