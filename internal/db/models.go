package db

import "go.mongodb.org/mongo-driver/v2/bson"

type Event struct {
	ID   bson.ObjectID `bson:"_id"`
	Name string        `bson:"name"`
}
