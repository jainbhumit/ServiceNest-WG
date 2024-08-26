package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var mockClient *mongo.Client

func MockConnect() *mongo.Client {
	if mockClient == nil {
		uri := "mongodb://localhost:27017/"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		mockClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatal(err)
		}

		// Check the connection
		err = mockClient.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connected to MongoDB!")
	}
	return mockClient
}

func MockGetCollection(dbName, collectionName string) *mongo.Collection {
	return mockClient.Database(dbName).Collection(collectionName)
}

func MockDisconnect() {
	if mockClient != nil {
		err := mockClient.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Disconnected from MongoDB.")
	}
}
