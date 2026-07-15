#!/bin/sh
set -eu

PREFIX="${PREFIX:-$HOME/.local}"
BIN="$PREFIX/bin"
APPDIR="$PREFIX/share/kingo-linux-vpn"
DESKTOPDIR="$HOME/.config/autostart"
SYSTEMDDIR="$HOME/.config/systemd/user"

mkdir -p "$BIN" "$APPDIR" "$DESKTOPDIR" "$SYSTEMDDIR"

if [ -f "./build/kingo-linux-vpn-daemon" ]; then
    install -m 755 ./build/kingo-linux-vpn-daemon "$BIN/kingo-linux-vpn-daemon"
fi

if [ -f "./build/kingo-linux-vpn-ui" ]; then
    install -m 755 ./build/kingo-linux-vpn-ui "$BIN/kingo-linux-vpn-ui"
fi

install -m 644 packaging/appimage/kingo-linux-vpn.desktop "$APPDIR/kingo-linux-vpn.desktop"
install -m 644 packaging/desktop/kingo-linux-vpn-autostart.desktop "$DESKTOPDIR/kingo-linux-vpn-autostart.desktop"
install -m 644 packaging/systemd/user/kingo-linux-vpn-daemon.service "$SYSTEMDDIR/kingo-linux-vpn-daemon.service"
install -m 644 packaging/systemd/user/kingo-linux-vpn-ui.service "$SYSTEMDDIR/kingo-linux-vpn-ui.service"

echo "Installed into $PREFIX"
echo "Enable services with:"
echo "  systemctl --user daemon-reload"
echo "  systemctl --user enable --now kingo-linux-vpn-daemon.service"
echo "  systemctl --user enable --now kingo-linux-vpn-ui.service"
