package db

import (
	// db "bdobot/db/connect"
	// b "bdobot/bdoapi"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertUser(u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := usersCollection.InsertOne(ctx, u)
	return err
}

func AddItemToSpecItems(userID int64, newItem ItemSpec) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": userID}
	update := bson.M{
		"$push": bson.M{"itemsOnSpec": newItem},
	}

	_, err := usersCollection.UpdateOne(ctx, filter, update)
	return err
}
