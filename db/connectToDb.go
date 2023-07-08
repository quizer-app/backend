package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func ConnectToDb() {
	uri := os.Getenv("MONGODB_URI")
	// production := os.Getenv("APP_PROD") == "true"

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic("Failed to connect to db")
	}

	db = client.Database("quizer")
}

func DisconnectFromDb() {
	db.Client().Disconnect(context.Background())
}

func GetCollection(collectionName string) *mongo.Collection {
	return db.Collection(collectionName)
}
