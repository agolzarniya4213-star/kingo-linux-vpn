
```markdown
# Kingo Linux VPN

A modern, secure, and high-performance VPN client for Linux, built with Go and Qt6/QML. It leverages the `sing-box` core to support modern protocols like VLESS, VMess, and Trojan.

## ✨ Features

- **Modern UI:** Built with Qt6/QML, featuring a dark theme and system tray integration.
- **High-Performance Core:** Powered by `sing-box` for optimal speed and reliability.
- **Smart Auto-Connect:** Automatically tests server latency and connects to the fastest available server.
- **Real-time Traffic Monitoring:** Live upload/download speed display via Clash API.
- **Subscription Management:** Fetches and parses V2Ray standard subscription links (Base64).
- **System Integration:** Runs as a secure background service using `systemd` with restricted capabilities.
- **Local Storage:** SQLite database for persisting server configurations.

## 🏗 Architecture

The project follows a strict Client-Daemon architecture:

- **Backend (Go):** Runs as a `systemd` service (`kingo-daemon`). Manages the `sing-box` process, handles subscription fetching, latency testing, and stores data in SQLite. Communicates with the UI via Unix Domain Sockets (UDS) using JSON.
- **Frontend (C++/Qt6):** Runs in user space. Provides a rich GUI, handles system tray interactions, and sends commands to the backend via UDS.

## 🚀 Building from Source

### Prerequisites
- Go 1.21+
- Qt6 (Quick, QuickControls2, Gui, Network, Widgets)
- CMake 3.16+
- `sing-box` installed and in system PATH

### Build Instructions
```bash
# Build Go Backend
go build -o build/kingo-linux-vpn-daemon ./cmd/kingo-linux-vpn/

# Build Qt UI
mkdir -p build/qt && cd build/qt
cmake ../../qt -DCMAKE_BUILD_TYPE=Release
cmake --build .
cd ../..
```

## 📦 Installation

To install the application system-wide:
```bash
sudo ./scripts/install.sh
```
This will install the binaries, setup the `systemd` service, and add the application to your desktop menu.

## 🛡 Security

The daemon runs with strict `systemd` hardening:
- `ProtectSystem=strict`
- `ProtectHome=true`
- `NoNewPrivileges=true`
- Only requires `CAP_NET_ADMIN` and `CAP_NET_BIND_SERVICE` capabilities for TUN mode and port binding.

## 📄 License

This project is licensed under the MIT License.
```

