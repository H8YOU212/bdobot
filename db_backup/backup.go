package backup

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Backup() {
	_ = godotenv.Load("./.env")
	uri := os.Getenv("uri")

	ctx := context.Background()

	// Подключение к исходной БД
	sourceClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	sourceDB := sourceClient.Database("bdobot")

	// Подключение к целевой БД
	targetClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	targetDB := targetClient.Database("bdobot_backup")

	collections, _ := sourceDB.ListCollectionNames(ctx, struct{}{})
	for _, collName := range collections {
		sourceColl := sourceDB.Collection(collName)
		targetColl := targetDB.Collection(collName)

		cursor, _ := sourceColl.Find(ctx, struct{}{})
		var docs []interface{}
		cursor.All(ctx, &docs)
		if len(docs) > 0 {
			targetColl.InsertMany(ctx, docs)
			fmt.Println("Скопирована коллекция:", collName)
		}
	}
}