package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	fuzz "github.com/google/gofuzz" // Used for creating random todo items
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-change-stream-demo/internal/app"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://mongo1:30001")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Get a handle for your collection
	collection := client.Database("test").Collection("event")

	/*
		_, err = collection.Indexes().DropAll(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	*/

	// Create Index Model for the TTL index
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(10),
	}

	// Request the server to create the TTL index. If the index already exists the server simply ignores
	// the request. However, not if the options of the given Index Model change this will lead to an array. It is then
	// required to drop the index first.
	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	// Insert random event items at two second intervals
	insert(collection)
}

func randomInt(upperLimit int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(upperLimit)))
	if err != nil {
		panic(err)
	}

	return int(nBig.Int64())
}

func insert(collection *mongo.Collection) {
	unicodeRanges := fuzz.UnicodeRanges{
		{First: 'a', Last: 'z'},
		{First: '0', Last: '9'},
	}

	f := fuzz.New().Funcs(unicodeRanges.CustomStringFuzzFunc())

	for {
		t := app.Event{}

		t.MessageID = [2]string{"d692e245-a93d-45b7-910f-e61d8b4f6035", "cce4011f-438c-40b9-befe-993350d88808"}[randomInt(2)] + ":" + uuid.New().String()
		f.Fuzz(&t.Message)
		t.CreatedAt = time.Now().UTC()

		fmt.Println(t.MessageID)

		_, err := collection.InsertOne(context.TODO(), t)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)
	}
}
