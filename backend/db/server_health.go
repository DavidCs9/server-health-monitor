// db/server_health.go

package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertServerStatus inserts a new server status document into the database
func InsertServerStatus(serverStatus ServerStatus) error {
	collection := GetCollection()

	// Insert the server status
	_, err := collection.InsertOne(context.Background(), serverStatus)
	if err != nil {
		return err
	}

	return nil
}

// GetServerHealth retrieves server health data for a specific server within a time range
func GetServerHealth(serverURL string) ([]ServerStatus, error) {
	// Get the collection (assuming this is your MongoDB collection)
	collection := GetCollection()
	log.Printf("Getting server health for %s", serverURL)

	// Construct the filter
	filter := bson.M{
		"server_url": serverURL,
	}

	sort := bson.M{
		"timestamp": -1,
	}

	// Execute the query
	cursor, err := collection.Find(context.Background(), filter, &options.FindOptions{
		Sort: sort})
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Prepare to store results
	var results []ServerStatus

	// Process the results
	if err := cursor.All(context.Background(), &results); err != nil {
		log.Printf("Error reading results: %v", err)
		return nil, err
	}

	return results, nil
}
