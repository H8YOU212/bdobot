package db

import (
	// db "bdobot/db/connect"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// var client *mongo.Client
var result bool
var user bson.M

func SearchByID(database, collection string, filter int64) (bool, error) {
	Conn()
	if client == nil {
		return false, fmt.Errorf("client uninit: %v", client)
	}
	defer Dconn()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user = bson.M{"id": fmt.Sprintf("%v", filter)}

	db := client.Database(database)
	db.Collection(collection).FindOne(ctx, user)
	
	return result, nil
}
