package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUserItems(id int, newItems []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}
	update := bson.M{
		"$set": bson.M{"itemsOnSpec": newItems},
	}

	_, err := usersCollection.UpdateOne(ctx, filter, update)
	return err
}
