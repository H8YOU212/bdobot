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
var result bool
var user bson.M

func FindUserByID(id int64) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"id": id}

	var user User
	err := usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserExists(userID int64) (bool, error) {
	coll, err := GetUsersCollection()
	if err != nil {
		return false, err
	}

	filter := bson.M{"id": userID}
	var user User
	err = coll.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetUsersCollection() (*mongo.Collection, error) {
	if usersCollection == nil {
		return nil, fmt.Errorf("коллекция пользователей не инициализирована")
	}
	return usersCollection, nil
}


func GetAllUsers() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	cursor , err := usersCollection.Find(ctx, bson.M{})
	if err != nil{
		return nil, err	
	}
	defer cursor.Close(ctx)
	
	var users []User

	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil{
			return nil, err
		}
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil{
		return nil, err
	}
	
	return users, nil
}
