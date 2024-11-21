package main

import (
	"encoding/json"
	"log"
	"net/http"
	health "server-health-monitor/pkg"
	"sync"
	"time"

	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers []string `yaml:"servers"`
}

var (
	serverStatus []health.ServerStatus
	mu           sync.Mutex
)

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

	var statuses []health.ServerStatus
	for status := range statusChannel {
		statuses = append(statuses, status)
	}

	mu.Lock()
	serverStatus = statuses
	mu.Unlock()

}

func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	if len(serverStatus) == 0 {
		http.Error(w, "No health data available", http.StatusServiceUnavailable)
		return
	}

	err := json.NewEncoder(w).Encode(serverStatus)
	if err != nil {
		http.Error(w, "Failed to encode serverStatus", http.StatusInternalServerError)
		log.Printf("Failed to encode serverStatus: %v", err)
		return
	}
}

func getHome(w http.ResponseWriter, r *http.Request) {
	// Set Content-Type header to text/html
	w.Header().Set("Content-Type", "text/html")

	// Define the HTML response
	htmlResponse := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Server Status</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #e0f7fa;
				color: #333;
				text-align: center;
				padding: 50px;
			}
			.status-box {
				display: inline-block;
				padding: 20px;
				border: 2px dashed #66bb6a;
				border-radius: 10px;
				background-color: #ffffff;
			}
			h1 {
				color: #66bb6a;
			}
			p {
				font-size: 1.2em;
			}
		</style>
	</head>
	<body>
		<div class="status-box">
			<h1>The Server is Up and Running! 🚀</h1>
			<p>Everything is working as expected. 🎉</p>
		</div>
	</body>
	</html>
	`

	// Write the HTML response to the ResponseWriter
	_, err := w.Write([]byte(htmlResponse))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func main() {
	config := loadConfig("servers.yaml")

	// Perform an initial health check
	monitorServers(config.Servers)

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			monitorServers(config.Servers)
		}
	}()

	// setup http routes using gorilla/mux
	r := mux.NewRouter()
	r.HandleFunc("/health", getHealthHandler).Methods("GET")
	r.HandleFunc("/", getHome)

	// start http server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
