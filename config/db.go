package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

// dbConnect initializes the MongoDB client connection
func DbConnect() error {
	var err error
	uri := "mongodb+srv://deep82500:deep82500@deep.jqe1i.mongodb.net/?retryWrites=true&w=majority&appName=deep"

	clientAddress := options.Client().ApplyURI(uri)
	client, err = mongo.Connect(context.TODO(), clientAddress)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return err
	}

	// Optionally check if the connection is alive
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Error pinging MongoDB:", err)
		return err
	}

	return nil
}

// GetCollection
func GetCollection() *mongo.Collection {
	if client == nil {
		log.Fatal("MongoDB client is not initialized")
	}
	collection = client.Database("new_employee").Collection("employee_list")

	return collection
}
