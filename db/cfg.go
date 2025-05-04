package db

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	// b "bdobot/bdoapi"
	// "go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID          int64      `bson:"id"`
	Name        string     `bson:"name"`
	ItemsOnSpec []ItemSpec `bson:"itemsOnSpec"`
}

type ItemSpec struct {
	ID              int    `bson:"id"`
	Name            string `bson:"name"`
	Price           int    `bson:"price"`
	ItemStartPrice  int    `bson:"itemStartPrice"`
	ItemTargetPrice int    `bson:"itemEndPrice"`
}

var (
	client          *mongo.Client
	usersCollection *mongo.Collection
	connectOnce     sync.Once
	connectErr      error
)
