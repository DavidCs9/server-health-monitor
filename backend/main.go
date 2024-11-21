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

	log.Printf("ServerStatus before encoding: %+v", serverStatus)

	err := json.NewEncoder(w).Encode(serverStatus)
	if err != nil {
		http.Error(w, "Failed to encode serverStatus", http.StatusInternalServerError)
		log.Printf("Failed to encode serverStatus: %v", err)
		return
	}

	log.Println("Successfully responded with serverStatus.")
}

func getHome(w http.ResponseWriter, r *http.Request) {
	// Define response as a map
	response := map[string]bool{"success": true}

	w.Header().Set("Content-Type", "application/json")
	// Encode response as JSON with error handling
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
