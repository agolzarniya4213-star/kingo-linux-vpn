#!/bin/sh
set -e
DIR="$(cd "$(dirname "$0")" && pwd)"
"$DIR/build/kingo-linux-vpn-daemon" &
DAEMON_PID=$!
sleep 2
"$DIR/build/qt/kingo-linux-vpn-ui"
kill $DAEMON_PID 2>/dev/null
exit 0
