# Architecture

## Current state

The project now has:
- a CLI bootstrap
- a VPN engine runner
- Xray config generation
- a Unix socket daemon
- a client package for future UI integration

## Process model

1. `kingo-linux-vpn-daemon`
   - owns config directory
   - keeps server cache loaded
   - receives IPC commands over Unix socket
   - starts/stops the VPN engine

2. `kingo-linux-vpn`
   - CLI entrypoint for fetch, dry-run, and basic control
   - later becomes a thin client to the daemon

3. UI
   - talks to daemonclient
   - never manipulates engine directly

## Next milestone

- build a Qt/QML frontend
- add a server list view
- add connect/disconnect buttons
- add latency probe and favorites
