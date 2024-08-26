package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var client *mongo.Client

func Connect() *mongo.Client {
	if client == nil {
		uri := "mongodb://localhost:27017/"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var err error
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatal(err)
		}

		// Check the connection
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Connected to MongoDB!")
	}
	return client
}

func GetCollection(dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}

func Disconnect() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Disconnected from MongoDB.")
	}
}
