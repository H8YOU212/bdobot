package tgbot

import (
	"log"
	"os"

	//bdoapi "bdobot/bdoapi"
	h "bdobot/tgbot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func StartTelegramBotLoop() {
	err := godotenv.Load("tgbot/.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %v", err)
	}
	token := os.Getenv("TG_KEY")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка авторизации бота: %v", err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			h.HandleMessage(bot, update)
		} else if update.CallbackQuery != nil {
			// add "go" keyword if you want async handling
			// we're making chatState with mutex, so it **SHOULD** be good
			go h.HandleCallback(bot, update)
		}
	}
}
