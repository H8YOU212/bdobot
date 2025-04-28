package db

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

// var client *mongo.Client

func Delete() (*mongo.DeleteResult, error) {
	Conn()
	if client == nil {
		return nil, fmt.Errorf("client uninit: %v", client)
	}
	defer Dconn()
	
	
	
	return nil, nil
}
