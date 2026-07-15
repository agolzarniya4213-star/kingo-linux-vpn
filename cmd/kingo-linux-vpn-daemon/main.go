package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kingo-linux/vpn/internal/config"
	"github.com/kingo-linux/vpn/internal/daemon"
)

func main() {
	cfg := config.Default()
	settings := config.DefaultSettings()

	d := daemon.New(cfg, settings)

	fmt.Printf("daemon listening on %s\n", d.SocketPath())
	if err := d.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "daemon error: %v\n", err)
		os.Exit(1)
	}
}
