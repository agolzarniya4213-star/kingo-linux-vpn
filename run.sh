#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
echo "Starting Kingo VPN Daemon with Root privileges for TUN routing..."
sudo "$DIR/build/kingo-linux-vpn-daemon" &
DAEMON_PID=$!
sleep 2
"$DIR/build/qt/kingo-linux-vpn-ui"
sudo kill $DAEMON_PID 2>/dev/null
