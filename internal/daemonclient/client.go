package daemonclient

import (
	"encoding/json"
	"net"
	"time"

	"github.com/kingo-linux/vpn/internal/daemon"
)

func Call(socketPath string, cmd daemon.Command, timeout time.Duration) (daemon.Response, error) {
	conn, err := net.DialTimeout("unix", socketPath, timeout)
	if err != nil {
		return daemon.Response{}, err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	if err := json.NewEncoder(conn).Encode(cmd); err != nil {
		return daemon.Response{}, err
	}

	var resp daemon.Response
	if err := json.NewDecoder(conn).Decode(&resp); err != nil {
		return daemon.Response{}, err
	}
	return resp, nil
}
