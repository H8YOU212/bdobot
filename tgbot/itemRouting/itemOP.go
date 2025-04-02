package itemrouting

import (
	// h "bdobot/tgbot/handlers"
	"bdobot/bdoapi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// type Item h.Item

var ( 
	itemIndexMap = make(map[int64]int)
	itemCache = make(map[int64][]bdoapi.Item)
	categoryCashe = make(map[int64]struct{
		mainC int
		subC int
	})
)

func ItemOP(bot *tgbotapi.BotAPI, update tgbotapi.Update, curIndex int) {
	switch curIndex{
	case 0:
		// h.EditMessage()
	}
}