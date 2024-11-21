package main

import (
	"fmt"
	"log"
	health "server-health-monitor/pkg"
	"sync"
	"time"

	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers []string `yaml:"servers"`
}

func loadConfig(filePath string) Config {
	var config Config
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading servers config file: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %v", err)
	}

	return config
}

func monitorServers(servers []string) {
	statusChannel := make(chan health.ServerStatus, len(servers))
	var wg sync.WaitGroup

	for _, server := range servers {
		wg.Add(1)
		go func(server string) {
			defer wg.Done()
			status := health.CheckServerWithTimeout(server, 500*time.Millisecond)
			statusChannel <- status
		}(server)
	}

	go func() {
		wg.Wait()
		close(statusChannel)
	}()

	for status := range statusChannel {
		fmt.Printf("Server: %s, IsUp: %v, Latency: %v\n", status.URL, status.IsUp, status.Latency)
	}

}

func main() {
	config := loadConfig("servers.yaml")
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		fmt.Println("Starting health check...")
		monitorServers(config.Servers)
		fmt.Println("Health check completed.")
	}
}
