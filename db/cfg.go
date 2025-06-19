package db

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	// b "bdobot/bdoapi"
	// "go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID          int64      		`bson:"id"`
	Name        string     		`bson:"name"`
	ItemsOnSpec []ItemSpec 		`bson:"itemsOnSpec"`
}

type ItemSpec struct {
	ID              int    		`bson:"id"`
	SID				int			`bson:"sid"`	
	Name            string 		`bson:"name"`
	Price           int			`bson:"price"`
	ItemStartPrice  int			`bson:"itemStartPrice"`
	ItemTargetPrice int			`bson:"itemEndPrice"`
	Method 			int			`bson:"method"`
}

var (
	client          *mongo.Client
	usersCollection *mongo.Collection
	logCollection	*mongo.Collection
	connectOnce     sync.Once
	connectErr      error
)
