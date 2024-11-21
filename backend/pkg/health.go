package health

import (
	"context"
	"log"
	"net/http"
	"time"
)

type ServerStatus struct {
	URL     string `json:"url"`
	Latency string `json:"latency"` // Convert time.Duration to string for JSON
	IsUp    bool   `json:"is_up"`
}

func CheckServer(url string) ServerStatus {
	start := time.Now()
	resp, err := http.Get(url)
	latency := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
	}

	isUp := err == nil && resp.StatusCode == http.StatusOK

	if err != nil {
		log.Printf("Failed to check URL %s: %v", url, err)
	}

	return ServerStatus{
		URL:     url,
		Latency: latency.String(),
		IsUp:    isUp,
	}
}

func CheckServerWithTimeout(url string, timeout time.Duration) ServerStatus {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for URL %s: %v", url, err)
		return ServerStatus{
			URL:     url,
			Latency: "0s",
			IsUp:    false,
		}
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	latency := time.Since(start)

	if resp != nil {
		defer resp.Body.Close()
	}

	isUp := err == nil && resp.StatusCode == http.StatusOK

	if err != nil {
		log.Printf("Failed to check URL %s: %v", url, err)
	}

	return ServerStatus{
		URL:     url,
		Latency: latency.String(),
		IsUp:    isUp,
	}
}
