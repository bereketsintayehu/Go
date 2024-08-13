package bootstrap

import (
	"context"
	"log"
	"github.com/joho/godotenv"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)	

var client *mongo.Client

func ConnectDB() (*mongo.Client, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Connected to MongoDB!")
	return client, nil
}


func GetCollecton(collectionName string) *mongo.Collection {
	client, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database(os.Getenv("DB_NAME")).Collection(collectionName)
	return collection
}