package local

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"
)

// detectPublicIP fetches the public IP address from external services
// Returns empty string on failure or timeout
func detectPublicIP() string {
	// Try multiple services with timeout for redundancy
	services := []string{
		"https://api.ipify.org",
		"https://icanhazip.com",
		"https://ifconfig.me/ip",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for _, service := range services {
		req, err := http.NewRequestWithContext(ctx, "GET", service, nil)
		if err != nil {
			continue
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			continue
		}

		if body, err := io.ReadAll(resp.Body); err == nil {
			resp.Body.Close()
			ip := strings.TrimSpace(string(body))
			if ip != "" {
				return ip
			}
		} else {
			resp.Body.Close()
		}
	}

	return ""
}
