package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client



func Conn() error {
	var err error
	err = godotenv.Load("./.env")
	if err != nil{
		return fmt.Errorf("err load env")
	}
	uri := os.Getenv("uri")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil{
		return fmt.Errorf("ошибка подключения: %v", err)
	}
	
	fmt.Printf("Успешное подключение: %v", client)

	return nil
}

func Dconn() error {
	if client == nil{
		return fmt.Errorf("client isn`t init")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil{
		return fmt.Errorf("Ошибка отключения: %v", err)
	}
	fmt.Println("Успешное отключение")
	return nil
}
