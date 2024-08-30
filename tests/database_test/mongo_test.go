package database_test

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"serviceNest/database"
	"testing"
	"time"
)

const (
	testDBName         = "testdb"
	testCollectionName = "testcollection"
)

func setup() (*mongo.Client, error) {
	client := database.Connect()                            // Connect to MongoDB
	err := client.Database(testDBName).Drop(context.TODO()) // Clean up the test database
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestConnect(t *testing.T) {
	client := database.Connect()
	if client == nil {
		t.Fatal("Expected MongoDB client, got nil")
	}

	// Ping to ensure the connection is successful
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}
}

func TestGetCollection(t *testing.T) {
	_, err := setup()
	if err != nil {
		t.Fatalf("Failed to setup MongoDB connection: %v", err)
	}
	defer database.Disconnect() // Ensure disconnection after test

	collection := database.GetCollection(testDBName, testCollectionName)
	if collection == nil {
		t.Fatal("Expected MongoDB collection, got nil")
	}

	// Verify that the collection can be used
	doc := bson.D{{"name", "test"}}
	_, err = collection.InsertOne(context.TODO(), doc)
	if err != nil {
		t.Fatalf("Failed to insert document into collection: %v", err)
	}

	var result bson.D
	err = collection.FindOne(context.TODO(), bson.D{{"name", "test"}}).Decode(&result)
	if err != nil {
		t.Fatalf("Failed to find document in collection: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Expected document to be found in collection")
	}
}
