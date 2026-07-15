# Packaging and startup

The current packaging layer includes:
- systemd user service for the daemon
- optional systemd user service for the UI
- autostart desktop entry
- AppImage launcher metadata
- install script for local testing

Expected runtime layout:
- `~/.local/bin/kingo-linux-vpn-daemon`
- `~/.local/bin/kingo-linux-vpn-ui`
- config under the app-specific config directory
