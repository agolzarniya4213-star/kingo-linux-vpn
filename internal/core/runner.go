package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/kingo-linux/vpn/internal/config"
	"github.com/kingo-linux/vpn/internal/model"
)

type State string

const (
	StateStopped  State = "stopped"
	StateStarting State = "starting"
	StateRunning  State = "running"
	StateFailed   State = "failed"
)

type Runner struct {
	mu        sync.Mutex
	state     State
	server    model.Server
	cancel    context.CancelFunc
	cmd       *exec.Cmd
	configDir string
	binary    string
	settings  config.Settings
	lastErr   error
}

func NewRunner(binaryPath string, settings config.Settings) *Runner {
	return &Runner{
		state:     StateStopped,
		configDir: filepath.Join(os.TempDir(), "kingo-linux-vpn"),
		binary:    binaryPath,
		settings:  settings,
	}
}

func (r *Runner) State() State {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.state
}

func (r *Runner) LastError() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.lastErr
}

func (r *Runner) CurrentServer() model.Server {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.server
}

func (r *Runner) Start(ctx context.Context, server model.Server) error {
	r.mu.Lock()
	if r.state == StateRunning || r.state == StateStarting {
		r.mu.Unlock()
		return errors.New("engine already running")
	}
	r.state = StateStarting
	r.server = server
	r.lastErr = nil
	r.mu.Unlock()

	cfgBytes, err := BuildXrayConfig(server, r.settings)
	if err != nil {
		r.setFailed(err)
		return err
	}

	cfgDir := filepath.Join(r.configDir, "runtime")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		r.setFailed(err)
		return err
	}
	cfgPath := filepath.Join(cfgDir, "xray.json")
	if err := os.WriteFile(cfgPath, cfgBytes, 0o644); err != nil {
		r.setFailed(err)
		return err
	}

	if _, err := os.Stat(r.binary); err != nil {
		r.setFailed(fmt.Errorf("engine binary not found at %s: %w", r.binary, err))
		return err
	}

	runCtx, cancel := context.WithCancel(ctx)
	cmd := exec.CommandContext(runCtx, r.binary, "run", "-c", cfgPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		cancel()
		r.setFailed(err)
		return err
	}

	r.mu.Lock()
	r.cancel = cancel
	r.cmd = cmd
	r.state = StateRunning
	r.mu.Unlock()

	go func() {
		_ = cmd.Wait()
		r.mu.Lock()
		defer r.mu.Unlock()
		if r.state == StateRunning {
			r.state = StateStopped
		}
		r.cmd = nil
		r.cancel = nil
	}()

	return nil
}

func (r *Runner) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state == StateStopped {
		return nil
	}
	if r.cancel != nil {
		r.cancel()
	}
	if r.cmd != nil && r.cmd.Process != nil {
		_ = r.cmd.Process.Kill()
		_, _ = r.cmd.Process.Wait()
	}
	r.cmd = nil
	r.cancel = nil
	r.state = StateStopped
	return nil
}

func (r *Runner) setFailed(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.state = StateFailed
	r.lastErr = err
}
