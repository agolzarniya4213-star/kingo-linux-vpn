package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/kingo-linux/vpn/internal/config"
	"github.com/kingo-linux/vpn/internal/daemon"
	"github.com/kingo-linux/vpn/internal/daemonclient"
)

type pageData struct {
	Name       string
	Version    string
	SocketPath string
	AppDir     string
}

var homePage = template.Must(template.New("home").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Name}}</title>
  <style>
    :root { color-scheme: dark; }
    body { font-family: system-ui, sans-serif; margin: 0; background: #0b0f14; color: #e8eef6; }
    header { padding: 20px 24px; border-bottom: 1px solid #223042; background: #0f1620; }
    main { padding: 24px; max-width: 1100px; margin: 0 auto; display: grid; gap: 16px; }
    .grid { display: grid; grid-template-columns: 1.1fr 0.9fr; gap: 16px; }
    .card { background: #121a25; border: 1px solid #223042; border-radius: 18px; padding: 18px; box-shadow: 0 8px 24px rgba(0,0,0,.24); }
    .muted { color: #9db0c7; }
    .row { display: flex; gap: 10px; flex-wrap: wrap; align-items: center; }
    button { border: 0; border-radius: 12px; padding: 10px 14px; background: #2f7cf6; color: white; font-weight: 600; cursor: pointer; }
    button.secondary { background: #233244; }
    button.danger { background: #c94f4f; }
    textarea { width: 100%; box-sizing: border-box; border-radius: 12px; border: 1px solid #2b3b4f; background: #0f1620; color: #e8eef6; padding: 10px 12px; min-height: 96px; resize: vertical; }
    .servers { display: grid; gap: 10px; }
    .server { border: 1px solid #223042; border-radius: 14px; padding: 12px; background: #0f1620; }
    .server strong { display: block; margin-bottom: 4px; }
    .pill { display: inline-block; font-size: 12px; padding: 4px 10px; border-radius: 999px; background: #1e2a39; color: #bcd0e3; }
    pre { white-space: pre-wrap; word-break: break-word; background: #0f1620; padding: 12px; border-radius: 12px; border: 1px solid #223042; overflow-x: auto; }
    .status { font-size: 18px; font-weight: 700; }
  </style>
</head>
<body>
  <header>
    <div style="max-width:1100px;margin:0 auto;">
      <div class="row" style="justify-content:space-between;">
        <div>
          <div class="status" id="statusText">Loading…</div>
          <div class="muted">{{.Name}} {{.Version}} · Socket: {{.SocketPath}}</div>
        </div>
        <div class="row">
          <button onclick="refreshServers()">Refresh</button>
          <button class="secondary" onclick="reloadStatus()">Status</button>
          <button class="danger" onclick="stopEngine()">Stop</button>
        </div>
      </div>
    </div>
  </header>

  <main>
    <div class="grid">
      <section class="card">
        <h2>Servers</h2>
        <div class="muted">Pulled from the daemon cache.</div>
        <div id="servers" class="servers" style="margin-top:14px;"></div>
      </section>

      <section class="card">
        <h2>Connect</h2>
        <div class="muted">Paste a config or start a selected server.</div>
        <div style="margin-top:12px;">
          <label class="muted">Server JSON (optional)</label>
          <textarea id="serverJson" placeholder='{"name":"My Server","config":"vless://...","protocol":"VLESS"}'></textarea>
        </div>
        <div class="row" style="margin-top:12px;">
          <button onclick="startEngine()">Start</button>
          <button class="secondary" onclick="reloadStatus()">Reload status</button>
        </div>
      </section>
    </div>

    <section class="card">
      <h2>Raw status</h2>
      <pre id="statusJson">…</pre>
    </section>
  </main>

<script>
async function api(path, payload) {
  const res = await fetch(path, {
    method: payload ? 'POST' : 'GET',
    headers: payload ? {'Content-Type': 'application/json'} : undefined,
    body: payload ? JSON.stringify(payload) : undefined
  });
  return await res.json();
}

async function reloadStatus() {
  const data = await api('/api/status');
  document.getElementById('statusJson').textContent = JSON.stringify(data, null, 2);
  document.getElementById('statusText').textContent = data.ok ? ('State: ' + (data.state || 'unknown')) : ('Error: ' + (data.error || 'unknown'));
  renderServers(data.servers || []);
}

function renderServers(list) {
  const root = document.getElementById('servers');
  root.innerHTML = '';
  if (!list.length) {
    root.innerHTML = '<div class="muted">No servers cached yet. Use Refresh.</div>';
    return;
  }
  for (const s of list) {
    const el = document.createElement('div');
    el.className = 'server';

    const strong = document.createElement('strong');
    strong.textContent = s.name || 'Unnamed';
    el.appendChild(strong);

    const info = document.createElement('div');
    info.className = 'muted';
    const ping = (s.ping_ms !== undefined && s.ping_ms !== null) ? s.ping_ms : -1;
    info.textContent = (s.protocol || 'Unknown') + ' · ping ' + ping + ' ms';
    el.appendChild(info);

    const row = document.createElement('div');
    row.className = 'row';
    row.style.marginTop = '10px';

    const pill = document.createElement('span');
    pill.className = 'pill';
    pill.textContent = s.group || 'servers';
    row.appendChild(pill);

    const btn = document.createElement('button');
    btn.className = 'secondary';
    btn.textContent = 'Connect';
    btn.onclick = function() { startFrom(s); };
    row.appendChild(btn);

    el.appendChild(row);
    root.appendChild(el);
  }
}

async function refreshServers() {
  await api('/api/refresh', {});
  await reloadStatus();
}

async function stopEngine() {
  await api('/api/stop', {});
  await reloadStatus();
}

async function startEngine() {
  const raw = document.getElementById('serverJson').value.trim();
  let payload = {};
  if (raw) {
    try { payload = JSON.parse(raw); } catch (e) { alert('Invalid JSON'); return; }
  }
  await api('/api/start', payload);
  await reloadStatus();
}

async function startFrom(server) {
  document.getElementById('serverJson').value = JSON.stringify(server, null, 2);
  await startEngine();
}

reloadStatus();
</script>
</body>
</html>`))

func main() {
	cfg := config.Default()
	socketPath := cfg.AppDir + "/run/daemon.sock"

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = homePage.Execute(w, pageData{
			Name:       "Kingo Linux VPN",
			Version:    "0.1.0",
			SocketPath: socketPath,
			AppDir:     cfg.AppDir,
		})
	})

	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, status(socketPath))
	})
	mux.HandleFunc("/api/refresh", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, call(socketPath, daemon.Command{Action: "refresh"}))
	})
	mux.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, call(socketPath, daemon.Command{Action: "list"}))
	})
	mux.HandleFunc("/api/start", func(w http.ResponseWriter, r *http.Request) {
		var cmd daemon.Command
		if r.Method == http.MethodPost {
			_ = json.NewDecoder(r.Body).Decode(&cmd)
		}
		cmd.Action = "start"
		writeJSON(w, call(socketPath, cmd))
	})
	mux.HandleFunc("/api/stop", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, call(socketPath, daemon.Command{Action: "stop"}))
	})

	srv := &http.Server{
		Addr:              "127.0.0.1:8090",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("Kingo Linux VPN UI at http://%s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func call(socketPath string, cmd daemon.Command) any {
	resp, err := daemonclient.Call(socketPath, cmd, 5*time.Second)
	if err != nil {
		return map[string]any{"ok": false, "error": err.Error()}
	}
	return resp
}

func status(socketPath string) any {
	resp, err := daemonclient.Call(socketPath, daemon.Command{Action: "status"}, 5*time.Second)
	if err != nil {
		return map[string]any{"ok": false, "error": err.Error()}
	}
	return resp
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
