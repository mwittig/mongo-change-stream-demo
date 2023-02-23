package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"

	fuzz "github.com/google/gofuzz" // Used for creating random todo items
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type todo struct {
	Item string `bson:"item"`
	Done bool   `bson:"done"`
}

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
	FullDocument  todo                `bson:"fullDocument"`
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
	collection := client.Database("test").Collection("todo")

	// Watches the todo collection and prints out any changed documents
	go watch(collection)

	// Inserts random todo items at two second intervals
	insert(collection)

}

func watch(collection *mongo.Collection) {
	matchPipeline := bson.D{
		{
			"$match", bson.D{
				{"operationType", "insert"},
				{"fullDocument.item", primitive.Regex{Pattern: "^[a-z][a-z]", Options: ""}},
			},
		},
	}

	// Watch the todo collection
	cs, err := collection.Watch(context.TODO(), mongo.Pipeline{matchPipeline})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Whenever there is a new change event, decode the change event and print some information about it
	for cs.Next(context.TODO()) {
		var changeEvent changeEvent

		err := cs.Decode(&changeEvent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Change Event: %v\nTodo Item: %v\nDone: %v\n\n", changeEvent.OperationType, changeEvent.FullDocument.Item, changeEvent.FullDocument.Done)
	}
}

func insert(collection *mongo.Collection) {
	unicodeRanges := fuzz.UnicodeRanges{
		{'a', 'z'},
		{'0', '9'}, // You can also use 0x0030 as 0, 0x0039 as 9.
	}

	f := fuzz.New().Funcs(unicodeRanges.CustomStringFuzzFunc())

	for {
		t := todo{}

		f.Fuzz(&t.Item)
		f.Fuzz(&t.Done)

		_, err := collection.InsertOne(context.TODO(), t)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(2 * time.Second)
	}
}
