# Kingo Linux VPN

A Linux-native rewrite inspired by Kingo VPN's workflow:
- subscription refresh
- server list management
- latency testing
- connect / disconnect orchestration
- update and notification hooks

## Initial architecture

- `cmd/kingo-linux-vpn`: CLI entrypoint and future daemon bootstrap
- `internal/model`: shared data contracts
- `internal/config`: app configuration and paths
- `internal/storage`: local persistence
- `internal/fetcher`: subscription fetching and parsing
- `internal/core`: VPN engine orchestration
- `ui/`: Qt/QML desktop frontend (next step)

## Build

```bash
go build ./...
```

This repository is the first scaffold for the Linux version.

## Native UI

A Qt/QML scaffold now lives in `qt/` as the long-term desktop frontend.

## Packaging

Systemd user services, autostart entries, and AppImage metadata are now in `packaging/`.
