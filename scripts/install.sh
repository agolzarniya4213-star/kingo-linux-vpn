#!/bin/bash
set -e

# بررسی دسترسی root
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root (use sudo)"
  exit 1
fi

echo ">>> Installing Kingo Linux VPN Daemon..."
cp build/kingo-linux-vpn-daemon /usr/local/bin/
chmod +x /usr/local/bin/kingo-linux-vpn-daemon

echo ">>> Installing Systemd Service..."
cp deploy/kingo-daemon.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable kingo-daemon
systemctl restart kingo-daemon

echo ">>> Installing Kingo Linux VPN UI..."
cp build/qt/kingo-linux-vpn-ui /usr/local/bin/
chmod +x /usr/local/bin/kingo-linux-vpn-ui

echo ">>> Installing Desktop Entry..."
cp deploy/kingo-linux-vpn.desktop /usr/share/applications/

echo ">>> Installation Complete!"
echo "The daemon is running in the background via systemd."
echo "You can now launch 'Kingo Linux VPN' from your application menu."
