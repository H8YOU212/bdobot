package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() error {
	connectOnce.Do(func() {
		connectErr = Conn()
	})
	return connectErr
}

func Conn() error {
	// Загружаем переменные окружения
	err := godotenv.Load("db/.env")
	if err != nil {
		log.Println("Ошибка загрузки .env файла")
		return fmt.Errorf("не удалось загрузить .env файл: %v", err)
	}
	log.Println("Загрузка .env успешна")

	// Чтение URI подключения из переменной окружения
	uri := os.Getenv("uri")
	if uri == "" {
		log.Println("URI подключения пустое")
		return fmt.Errorf("не найдено значение для переменной окружения 'uri'")
	}
	// log.Printf("Строка подключения: %s\n", uri)

	// Подключаемся к MongoDB с тайм-аутом 10 секунд
	log.Println("Попытка подключиться к MongoDB...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v\n", err)
		return fmt.Errorf("ошибка подключения к MongoDB: %v", err)
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Ошибка пинга MongoDB: %v\n", err)
		return fmt.Errorf("ошибка пинга MongoDB: %v", err)
	}

	log.Println("Подключение к MongoDB успешно установлено")

	// Чтение имени базы данных и коллекции из переменных окружения
	db := os.Getenv("dbname")
	coll := os.Getenv("collname")
	logcoll := "logCollection"
	if db == "" || coll == "" {
		log.Println("Ошибка: не указаны имя базы данных или имя коллекции")
		return fmt.Errorf("не указаны имя базы данных или коллекции")
	}
	log.Printf("Используемая база данных: %s, коллекция: %s\n", db, coll)

	// Инициализация коллекции
	usersCollection = client.Database(db).Collection(coll)
	if usersCollection == nil {
		log.Println("Ошибка инициализации коллекции пользователей")
		return fmt.Errorf("не удалось инициализировать коллекцию %s в базе данных %s", coll, db)
	}

	logCollection = client.Database(db).Collection(logcoll)

	log.Println("Коллекция пользователей инициализирована успешно")
	return nil
}
