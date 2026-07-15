package daemon

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kingo-linux/vpn/internal/config"
	"github.com/kingo-linux/vpn/internal/core"
	"github.com/kingo-linux/vpn/internal/fetcher"
	"github.com/kingo-linux/vpn/internal/model"
	"github.com/kingo-linux/vpn/internal/storage"
)

type Command struct {
	Action string       `json:"action"`
	Server model.Server `json:"server,omitempty"`
}

type Response struct {
	OK      bool           `json:"ok"`
	Message string         `json:"message,omitempty"`
	State   string         `json:"state,omitempty"`
	Server  *model.Server  `json:"server,omitempty"`
	Err     string         `json:"error,omitempty"`
	Servers []model.Server `json:"servers,omitempty"`
}

type Daemon struct {
	cfg      config.Config
	settings config.Settings
	runner   *core.Runner

	mu    sync.Mutex
	store storage.Store
	ready bool
}

func New(cfg config.Config, settings config.Settings) *Daemon {
	return &Daemon{
		cfg:      cfg,
		settings: settings,
		runner:   core.NewRunner(cfg.EngineBinary, settings),
	}
}

func (d *Daemon) Run(ctx context.Context) error {
	if err := d.loadStore(); err != nil {
		return err
	}
	if err := os.MkdirAll(d.socketDir(), 0o755); err != nil {
		return err
	}
	_ = os.Remove(d.socketPath())

	ln, err := net.Listen("unix", d.socketPath())
	if err != nil {
		return err
	}
	defer ln.Close()

	go func() {
		<-ctx.Done()
		_ = ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			continue
		}
		go d.handleConn(conn)
	}
}

func (d *Daemon) handleConn(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(bufio.NewReader(conn))
	enc := json.NewEncoder(conn)

	var cmd Command
	if err := dec.Decode(&cmd); err != nil {
		_ = enc.Encode(Response{OK: false, Err: err.Error()})
		return
	}

	resp := d.execute(context.Background(), cmd)
	_ = enc.Encode(resp)
}

func (d *Daemon) execute(ctx context.Context, cmd Command) Response {
	switch strings.ToLower(strings.TrimSpace(cmd.Action)) {
	case "status":
		st := d.runner.State()
		cur := d.runner.CurrentServer()
		return Response{OK: true, State: string(st), Server: &cur, Servers: d.storeServers()}
	case "list":
		return Response{OK: true, State: string(d.runner.State()), Servers: d.storeServers()}
	case "start":
		srv := cmd.Server
		if err := srv.Validate(); err != nil {
			if len(d.storeServers()) > 0 {
				srv = d.storeServers()[0]
			}
		}
		if err := d.runner.Start(ctx, srv); err != nil {
			return Response{OK: false, Err: err.Error(), State: string(d.runner.State())}
		}
		return Response{OK: true, Message: "started", State: string(d.runner.State()), Server: &srv}
	case "stop":
		if err := d.runner.Stop(); err != nil {
			return Response{OK: false, Err: err.Error(), State: string(d.runner.State())}
		}
		return Response{OK: true, Message: "stopped", State: string(d.runner.State())}
	case "refresh":
		if err := d.refresh(); err != nil {
			return Response{OK: false, Err: err.Error()}
		}
		return Response{OK: true, Message: "refreshed", Servers: d.storeServers()}
	default:
		return Response{OK: false, Err: "unknown action"}
	}
}

func (d *Daemon) loadStore() error {
	s, err := storage.Load(d.cfg.ServersFile)
	if err != nil {
		return err
	}
	d.mu.Lock()
	d.store = s
	d.ready = true
	d.mu.Unlock()
	return nil
}

func (d *Daemon) refresh() error {
	result, err := fetchSubscription(d.cfg.SubscriptionURL, 20*time.Second)
	if err != nil {
		return err
	}
	d.mu.Lock()
	d.store.Servers = append(d.store.Servers, result...)
	s := d.store
	d.mu.Unlock()
	return storage.Save(d.cfg.ServersFile, s)
}

func fetchSubscription(url string, timeout time.Duration) ([]model.Server, error) {
	res, err := fetcher.FetchSubscription(url, timeout)
	if err != nil {
		return nil, err
	}
	return res.Servers, nil
}

func (d *Daemon) storeServers() []model.Server {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]model.Server, len(d.store.Servers))
	copy(out, d.store.Servers)
	return out
}

func (d *Daemon) socketDir() string {
	return filepath.Join(d.cfg.AppDir, "run")
}

func (d *Daemon) socketPath() string {
	return filepath.Join(d.socketDir(), "daemon.sock")
}

func (d *Daemon) Ready() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.ready
}

func ErrDaemonNotAvailable() error {
	return errors.New("daemon not available")
}

func (d *Daemon) SocketPath() string {
	return d.socketPath()
}

func (d *Daemon) Settings() config.Settings {
	return d.settings
}

func (d *Daemon) RunnerState() string {
	return string(d.runner.State())
}

func (d *Daemon) EngineBinary() string {
	return d.cfg.EngineBinary
}

func (d *Daemon) Shutdown() {
	_ = d.runner.Stop()
}

func (d *Daemon) Send(cmd Command) (Response, error) {
	conn, err := net.Dial("unix", d.socketPath())
	if err != nil {
		return Response{}, err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err := json.NewEncoder(conn).Encode(cmd); err != nil {
		return Response{}, err
	}
	var resp Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return Response{}, err
	}
	return resp, nil
}

func (d *Daemon) Status() (Response, error) {
	return d.Send(Command{Action: "status"})
}
