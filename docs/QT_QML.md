# Qt/QML milestone

The native desktop UI now has a working controller/model handshake:
- daemon responses update the server list model
- JSON input can be parsed and sent to `start`
- the UI can display daemon errors through `lastError`

The next step is packaging and tray integration.
