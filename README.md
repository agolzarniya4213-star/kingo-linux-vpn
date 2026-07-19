
<div align="center">

# 🛡️ Kingo Linux VPN

### A Modern, Secure, and High-Performance VPN Client for Linux

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)
![Qt](https://img.shields.io/badge/Qt-6-41CD52?logo=qt&logoColor=white)
![Sing-box](https://img.shields.io/badge/Core-Sing--box-blueviolet)
![License](https://img.shields.io/badge/License-MIT-green)

</div>

## 📖 Overview
Kingo Linux VPN is an enterprise-grade, lightweight VPN client designed specifically for Linux environments. Built with a strict **Client-Daemon architecture**, it separates the privileged network operations (Daemon) from the user interface (Client) to ensure maximum security and stability.

## ✨ Key Features
- 🚀 **High-Performance Core:** Powered by `sing-box` for optimal speed and modern protocol support (VLESS, VMess, Trojan).
- 🎨 **Modern QML UI:** Built with Qt6/QML, featuring a dark theme, responsive design, and System Tray integration.
- 🧠 **Smart Auto-Connect:** Automatically tests server latency and connects to the fastest available server.
- 📊 **Real-time Monitoring:** Live upload/download speed display via Clash API.
- 🔒 **Systemd Integration:** Runs as a hardened background service with strict capability bounding (`CAP_NET_ADMIN`).
- 💾 **Local Storage:** SQLite database for persisting server configurations and latency history.

## 🏗️ Architecture

```text
+-------------------------------------------------------------+
|                    User Space (Run as User)                 |
|                                                             |
|  +-------------------+              +--------------------+  |
|  |   Qt6/QML UI      |  <== UDS ==> |  Go IPC Server     |  |
|  | (C++ Bridge)      |  (JSON)      |  (Local Listener)  |  |
|  +-------------------+              +---------+----------+  |
|                                               |             |
|                                    +----------+----------+  |
|                                    |  Business Logic     |  |
|                                    +---------------------+  |
+-------------------------------------------------------------+
                                               |
                                     (Systemd / Polkit)
                                               |
+-------------------------------------------------------------+
|                Privileged Space (Run as Root)               |
|                                                             |
|  +-------------------+              +--------------------+  |
|  |  Go Daemon        | -----------> | Sing-box Engine    |  |
|  | (Process Manager) | <----------- | (TUN/Proxy Engine) |  |
|  +-------------------+              +--------------------+  |
+-------------------------------------------------------------+
```

## 🛠️ Building from Source

### Prerequisites
- Go 1.21+
- Qt6 (Quick, QuickControls2, Gui, Network, Widgets)
- CMake 3.16+
- `sing-box` binary installed and in system PATH

### Build Instructions
```bash
# 1. Build Go Backend
go build -o build/kingo-linux-vpn-daemon ./cmd/kingo-linux-vpn/

# 2. Build Qt UI
mkdir -p build/qt && cd build/qt
cmake ../../qt -DCMAKE_BUILD_TYPE=Release
cmake --build .
cd ../..

# 3. Run the application
./run.sh
```

## 📦 System Installation
To install the application system-wide (registers in application menu and starts daemon on boot):
```bash
sudo ./scripts/install.sh
```

## 🛡️ Security Hardening
The daemon runs with strict `systemd` hardening rules to minimize attack surface:
- `ProtectSystem=strict`: The daemon cannot modify system files.
- `ProtectHome=true`: User home directories are inaccessible.
- `NoNewPrivileges=true`: Prevents privilege escalation.
- `AmbientCapabilities=CAP_NET_ADMIN CAP_NET_BIND_SERVICE`: Only grants network administration rights, dropping all other root privileges.

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
