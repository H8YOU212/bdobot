package main

import (
	// "bdobot/tgbot/bdohandler"
	// "bdobot/bdoapi"
	"bdobot/db"
	"bdobot/tgbot"
	updateprices "bdobot/updatePrices"
	"log"
	// h "bdobot/tgbot/handlers"
)

func main() {
	// bdoapi.GetWorldMarketList(1, 2)
	// bdohandler.OutputItemData()

	if err := db.Conn(); err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	if _, err := db.GetUsersCollection(); err != nil {
		log.Fatalf("Ошибка подключения к коллекции пользователей: %v", err)
	}

	go func(){
		for{
			updateprices.StartUpdater()
		}
	}()

	tgbot.StartTelegramBotLoop()

	

	// h.FillItems(1, 1)
}
