package db

import (
	"context"
	"fmt"
	"queues-ocpi-poc/internal/configs"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type DB struct {
	Client *mongo.Client
}

func Connect() (*DB, error) {
	fmt.Println("Connecting to MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(configs.GetDBConnection()))
	if err != nil {
		fmt.Printf("error connecting to MongoDB: %v\n", err)
		return nil, err
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Printf("error pinging MongoDB: %v\n", err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB")

	return &DB{Client: client}, nil
}

func (c *DB) Disconnect(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}

func (c *DB) DropCollection(ctx context.Context, collection string) error {
	return c.Client.Database(configs.DB_NAME).Collection(collection).Drop(ctx)
}

func (c *DB) InsertDocument(ctx context.Context, collection string, document interface{}) error {
	collectionObject := c.Client.Database(configs.DB_NAME).Collection(collection)
	result, err := collectionObject.InsertOne(ctx, document)
	if err != nil {
		fmt.Printf("error inserting document: %v\n", err)
		return err
	}

	fmt.Println("Document inserted with ID: ", result.InsertedID)
	return nil
}

func (c *DB) ListDocuments(ctx context.Context, collection string) ([]Event, error) {
	fmt.Println("Listing documents from collection:", collection)
	collectionObject := c.Client.Database(configs.DB_NAME).Collection(collection)
	cursor, err := collectionObject.Find(ctx, bson.D{})
	if err != nil {
		fmt.Printf("error finding documents: %v\n", err)
		return nil, err
	}

	var documents []Event
	for cursor.Next(ctx) {
		fmt.Println("Decoding document...")
		var document Event
		err := cursor.Decode(&document)
		if err != nil {
			fmt.Printf("error decoding document: %v\n", err)
			return nil, err
		}
		fmt.Println("Document decoded: ", document)
		documents = append(documents, document)
	}
	return documents, nil
}
