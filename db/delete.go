package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

// var client *mongo.Client

func Delete(chatid int64, itemSpec ItemSpec) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": chatid}
	update := bson.M{
		"$pull": bson.M{
			"itemsOnSpec": bson.M{
				"id":				itemSpec.ID,
				"name":				itemSpec.Name,
				"price":			itemSpec.Price,
				"itemStartPrice":	itemSpec.ItemStartPrice,
				"itemEndPrice":		itemSpec.ItemTargetPrice,
				"method":			itemSpec.Method,
			},
		},
	}

	// op := usersCollection.FindOne(ctx, filter)
	res, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil{
		fmt.Printf("fail delete item, \n%v", err)
	}
	fmt.Printf("result is success, \n%v", res)

}
