// db/db.go

package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB client instance
var client *mongo.Client

// Connect initializes the MongoDB client
func Connect() error {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://davidcastro:Test123.@cluster0.ql7iv.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Check connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}

// GetCollection returns a reference to the server_status collection
func GetCollection() *mongo.Collection {
	return client.Database("server_monitor").Collection("server_status")
}
