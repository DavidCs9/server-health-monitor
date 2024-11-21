// db/models.go

package db

import "time"

// ServerStatus represents the health data of a server.
type ServerStatus struct {
	ServerURL string        `bson:"server_url" json:"server_url"`
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`
	Latency   time.Duration `bson:"latency" json:"latency"`
	IsUp      bool          `bson:"is_up" json:"is_up"`
}
