package health

import (
	"context"
	"net/http"
	"time"
)

type ServerStatus struct {
	URL     string
	Latency time.Duration
	IsUp    bool
}

func CheckServer(url string) ServerStatus {
	start := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(start)

	return ServerStatus{
		URL:     url,
		Latency: latency,
		IsUp:    err == nil && resp.StatusCode == http.StatusOK,
	}
}

func CheckServerWithTimeout(url string, timeout time.Duration) ServerStatus {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	latency := time.Since(start)
	isUp := err == nil && resp.StatusCode == http.StatusOK

	return ServerStatus{
		URL:     url,
		Latency: latency,
		IsUp:    isUp,
	}
}
