package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kingo-linux/vpn/internal/config"
	"github.com/kingo-linux/vpn/internal/core"
	"github.com/kingo-linux/vpn/internal/fetcher"
	"github.com/kingo-linux/vpn/internal/model"
	"github.com/kingo-linux/vpn/internal/storage"
	"github.com/kingo-linux/vpn/internal/version"
)

func main() {
	cfg := config.Default()
	settings := config.DefaultSettings()

	fetchCmd := flag.Bool("fetch", false, "fetch subscription and cache servers")
	dryRun := flag.Bool("dry-run", false, "build xray config and print it")
	startCmd := flag.Bool("start", false, "start the engine with the first saved server")
	statusCmd := flag.Bool("status", false, "print daemon/engine status")
	flag.Parse()

	fmt.Printf("%s %s\n", version.Name, version.Version)
	fmt.Printf("config dir: %s\n", cfg.AppDir)

	store, err := storage.Load(cfg.ServersFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load store: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("saved servers: %d\n", len(store.Servers))

	switch {
	case *fetchCmd:
		result, err := fetcher.FetchSubscription(cfg.SubscriptionURL, 20*time.Second)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch subscription: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("fetched servers: %d\n", len(result.Servers))
		store.Servers = append(store.Servers, result.Servers...)
		if err := storage.Save(cfg.ServersFile, store); err != nil {
			fmt.Fprintf(os.Stderr, "save store: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("subscription cached")
		return

	case *dryRun:
		srv := demoServer(store.Servers)
		cfgBytes, err := core.BuildXrayConfig(srv, settings)
		if err != nil {
			fmt.Fprintf(os.Stderr, "build config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(cfgBytes))
		return

	case *startCmd:
		srv := demoServer(store.Servers)
		runner := core.NewRunner(cfg.EngineBinary, settings)
		if err := runner.Start(context.Background(), srv); err != nil {
			fmt.Fprintf(os.Stderr, "start runner: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("engine started")
		return

	case *statusCmd:
		fmt.Println("status: local CLI mode; daemon socket will be added in the next step")
		return
	}

	fmt.Println("scaffold ready")
}

func demoServer(list []model.Server) model.Server {
	if len(list) > 0 {
		return list[0]
	}
	return model.Server{
		Name:     "placeholder",
		Config:   "vless://",
		Protocol: "VLESS",
		Group:    "servers",
		PingMS:   -1,
	}
}
