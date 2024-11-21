package health

import (
	"context"
	"log"
	"net/http"
	"server-health-monitor/db"
	"time"
)

func CheckServer(url string) db.ServerStatus {
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

	return db.ServerStatus{
		ServerURL: url,
		Latency:   latency,
		IsUp:      isUp,
		Timestamp: time.Now(),
	}
}

func CheckServerWithTimeout(url string, timeout time.Duration) db.ServerStatus {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request for URL %s: %v", url, err)
		return db.ServerStatus{
			ServerURL: url,
			Latency:   timeout,
			IsUp:      false,
			Timestamp: time.Now(),
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

	return db.ServerStatus{
		ServerURL: url,
		Latency:   latency,
		IsUp:      isUp,
		Timestamp: time.Now(),
	}
}
