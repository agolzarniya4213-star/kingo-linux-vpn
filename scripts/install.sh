#!/bin/bash
set -e

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (use sudo)"
  exit 1
fi

echo ">>> Installing Kingo Linux VPN Daemon..."
install -Dm755 build/kingo-linux-vpn-daemon /usr/local/bin/kingo-linux-vpn-daemon

echo ">>> Creating system user and directories..."
useradd -r -s /bin/false kingo-vpn || true
mkdir -p /var/lib/kingo-vpn /run/kingo-vpn /var/log/kingo-vpn
chown -R kingo-vpn:kingo-vpn /var/lib/kingo-vpn /run/kingo-vpn /var/log/kingo-vpn

echo ">>> Installing Systemd Service..."
install -Dm644 deploy/kingo-daemon.service /etc/systemd/system/kingo-daemon.service
systemctl daemon-reload
systemctl enable kingo-daemon
systemctl restart kingo-daemon

echo ">>> Installing Kingo Linux VPN UI..."
install -Dm755 build/qt/kingo-linux-vpn-ui /usr/local/bin/kingo-linux-vpn-ui

echo ">>> Installing Desktop Entry..."
install -Dm644 deploy/kingo-linux-vpn.desktop /usr/share/applications/kingo-linux-vpn.desktop

echo ">>> Installation Complete!"
