package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo-change-stream-demo/internal/app"
)

type documentKey struct {
	ID primitive.ObjectID `bson:"_id"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

// This is an example change event struct for inserts.
// It does not include all possible change event fields.
// You should consult the change event documentation for more info:
// https://docs.mongodb.com/manual/reference/change-events/
type changeEvent struct {
	ID            changeID            `bson:"_id"`
	OperationType string              `bson:"operationType"`
	ClusterTime   primitive.Timestamp `bson:"clusterTime"`
	FullDocument  app.Event           `bson:"fullDocument"`
	DocumentKey   documentKey         `bson:"documentKey"`
	Ns            namespace           `bson:"ns"`
}

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

	// Watches the event collection and prints out any changed documents. The pipeline for the server-side contains a
	// filter to let the server propagate change events for insert where MessageID start with the given UUID.
	watch(collection)
}

func watch(collection *mongo.Collection) {
	// Create a mact pipeline the server will process as part of change stream update. We are also concerned about
	// insert events for documents where the MessageID started with the given UUID string.
	matchPipeline := bson.D{
		{
			Key: "$match", Value: bson.D{
				{Key: "operationType", Value: "insert"},
				{Key: "fullDocument.message_id", Value: primitive.Regex{Pattern: "^cce4011f-438c-40b9-befe-993350d88808", Options: ""}},
			},
		},
	}

	// Watch the event collection
	cs, err := collection.Watch(context.TODO(), mongo.Pipeline{matchPipeline})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		fmt.Println(cs.Current)

		var eventMessage changeEvent

		err := cs.Decode(&eventMessage)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Change Event: %v\nMessageID: %v\nMessage: %v\nCluster Time: %s\n\n",
			eventMessage.OperationType, eventMessage.FullDocument.MessageID, eventMessage.FullDocument.Message, time.Unix(int64(eventMessage.ClusterTime.T), 0))
	}
}
