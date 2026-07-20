#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
echo "Starting Kingo VPN Daemon (Root required for TUN)..."
echo "Daemon logs are being saved to /tmp/kingo-daemon.log"
sudo "$DIR/build/kingo-linux-vpn-daemon" > /tmp/kingo-daemon.log 2>&1 &
DAEMON_PID=$!
sleep 2
"$DIR/build/qt/kingo-linux-vpn-ui"
sudo kill $DAEMON_PID 2>/dev/null
