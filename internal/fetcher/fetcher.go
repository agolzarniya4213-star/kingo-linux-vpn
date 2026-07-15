package fetcher

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kingo-linux/vpn/internal/model"
)

type Result struct {
	Servers []model.Server
	Raw     string
}

func FetchSubscription(subscriptionURL string, timeout time.Duration) (Result, error) {
	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(subscriptionURL)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	servers := make([]model.Server, 0)
	var raw strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		raw.WriteString(line)
		raw.WriteByte('\n')
		if strings.HasPrefix(line, "vless://") || strings.HasPrefix(line, "vmess://") || strings.HasPrefix(line, "trojan://") {
			servers = append(servers, model.ParseServerConfig(line, "servers"))
		}
	}
	if err := scanner.Err(); err != nil {
		return Result{}, err
	}

	return Result{Servers: servers, Raw: raw.String()}, nil
}
