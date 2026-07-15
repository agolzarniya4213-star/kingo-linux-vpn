# UI

The first Linux UI is a small local web panel that talks to the daemon.

Run:
- daemon: `go run ./cmd/kingo-linux-vpn-daemon`
- UI: `go run ./cmd/kingo-linux-vpn-ui`

This keeps the control plane stable before moving to Qt/QML.
