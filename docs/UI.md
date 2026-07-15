# UI milestone

The current UI milestone is a lightweight browser-based control panel that:
- reads daemon status
- lists cached servers
- starts/stops the engine
- triggers refresh

This is intentionally temporary and is the clean bridge toward the Qt/QML desktop UI. The UI contract is now stable:
- `status`
- `list`
- `start`
- `stop`
- `refresh`
