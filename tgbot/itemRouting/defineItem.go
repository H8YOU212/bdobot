package itemrouting

import (
	// h "bdobot/tgbot/handlers"
	"bdobot/bdoapi"
	"fmt"
	"log"

	// tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// type Item h.Item

var (
	itemIndexMap  = make(map[int64]int)
	itemCache     = make(map[int64][]bdoapi.Item)
	categoryCashe = make(map[int64]struct {
		mainC int
		subC  int
	})
)

func DefineItem(curIndex *int, mainC int, subC int) string {
	items, err := FillItems(mainC, subC)
	if err != nil {
		log.Println(err)
	}
	if *curIndex >= 0 && *curIndex < len(items) {
		item := items[*curIndex]
		message := fmt.Sprintf(fmt.Sprintf("Id предмета: %v, \nНазвание предмета: %v, \nЦена предмета: %v", item.ID, item.Name, item.Price))
		return message
	}
	return "errIndex"
}
