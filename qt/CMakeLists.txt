cmake_minimum_required(VERSION 3.21)
project(KingoLinuxVpnUI VERSION 0.1.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Quick QuickControls2 Network)

qt_standard_project_setup()

qt_add_executable(kingo-linux-vpn-ui
    src/main.cpp
    src/appcontroller.cpp
    src/appcontroller.h
    src/models/serverlistmodel.cpp
    src/models/serverlistmodel.h
)

qt_add_qml_module(kingo-linux-vpn-ui
    URI KingoVpn
    VERSION 1.0
    QML_FILES
        qml/Main.qml
        qml/components/ActionButton.qml
        qml/components/InfoCard.qml
        qml/components/ServerRow.qml
    RESOURCES
        assets/logo.svg
)

target_link_libraries(kingo-linux-vpn-ui
    PRIVATE Qt6::Quick Qt6::QuickControls2 Qt6::Network
)
// internal/daemon/tray_unix.go
package daemon

import (
    "context"
    "encoding/json"
    "fmt"
    "net"
    "os"
    "path/filepath"
)

type TrayStatus struct {
    Connected   bool   `json:"connected"`
    Connecting  bool   `json:"connecting"`
    ServerName  string `json:"server_name"`
    Error       string `json:"error,omitempty"`
    ServerList  []string `json:"server_list,omitempty"`
}

type TrayIPCServer struct {
    socketPath string
    statusChan chan TrayStatus
    commandChan chan string
}

func NewTrayIPCServer(configDir string) *TrayIPCServer {
    socketPath := filepath.Join(configDir, "tray.sock")
    return &TrayIPCServer{
        socketPath:  socketPath,
        statusChan:  make(chan TrayStatus, 10),
        commandChan: make(chan string, 10),
    }
}

func (s *TrayIPCServer) Start() error {
    // Remove old socket if exists
    os.Remove(s.socketPath)

    listener, err := net.Listen("unix", s.socketPath)
    if err != nil {
        return fmt.Errorf("failed to listen on socket: %w", err)
    }

    go func() {
        for {
            conn, err := listener.Accept()
            if err != nil {
                continue
            }

            go s.handleConnection(conn)
        }
    }()

    return nil
}

func (s *TrayIPCServer) handleConnection(conn net.Conn) {
    defer conn.Close()

    decoder := json.NewDecoder(conn)
    encoder := json.NewEncoder(conn)

    for {
        var command string
        if err := decoder.Decode(&command); err != nil {
            return
        }

        // Handle commands from tray
        switch command {
        case "status":
            status := <-s.statusChan
            encoder.Encode(status)
        case "connect":
            s.commandChan <- "connect"
        case "disconnect":
            s.commandChan <- "disconnect"
        case "quit":
            s.commandChan <- "quit"
        }
    }
}

func (s *TrayIPCServer) UpdateStatus(status TrayStatus) {
    s.statusChan <- status
}

func (s *TrayIPCServer) CommandChan() <-chan string {
    return s.commandChan
}
