#!/bin/bash
# scripts/setup-tray.sh - Setup and start tray icon

set -e

# Determine config directory
CONFIG_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/kingo-linux-vpn"
mkdir -p "$CONFIG_DIR"

# Create systemd user service directory
SYSTEMD_USER_DIR="${XDG_CONFIG_HOME:-$HOME/.config}/systemd/user"
mkdir -p "$SYSTEMD_USER_DIR"

# Copy service file
cp packaging/kingo-linux-vpn-tray.service "$SYSTEMD_USER_DIR/"

# Reload systemd daemon
systemctl --user daemon-reload

# Enable and start the service
systemctl --user enable kingo-linux-vpn-tray.service
systemctl --user start kingo-linux-vpn-tray.service

echo "Tray icon setup complete and started!"
echo "Service status: systemctl --user status kingo-linux-vpn-tray.service"
