package db

import (
	// db "bdobot/db/connect"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// var client *mongo.Client

func CollectionExists(database, collection string) (bool, error) {
	err := Conn()
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collNames, err := client.Database(database).ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return false, fmt.Errorf("Ошибка получения списка коллекций: %v", err)
	}

	for _, name := range collNames {
		if name == collection {
			return true, nil
		}
	}

	return false, nil
}

func Insert(database, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	err := Conn()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	coll := client.Database(database).Collection(collection)
	res, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	fmt.Println("Insert completes")
	return res, nil
}
