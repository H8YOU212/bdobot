package handlers

import (
	"bdobot/utils"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	ba "bdobot/bdoapi"
	"bdobot/tgbot/chatstate"
	"bdobot/tgbot/itemRouting"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Item struct {
	ID    int    `json:"id"`
	Sid   int    `json:"sid"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	// BasePrice    int            `json:"basePrice"`
	// History      map[string]int `json:"history"`
	// MainCategory int            `json:"mainCategory"`
	// SubCategory  int            `json:"subCategory"`
	// PriceMin     int            `json:"priceMin"`
	// PriceMax     int            `json:"priceMax"`
}

type Config struct {
	Token string
}

type Bot struct {
	bot      *tgbotapi.BotAPI
	handlers map[string]MessageHandler
}

type MessageHandler func(*tgbotapi.Message)

var ( 
	itemIndexMap = make(map[int64]int)
	itemCache = make(map[int64][]Item)
	categoryCashe = make(map[int64]struct{
		mainC int
		subC int
	})
)

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	commandName := update.Message.Command()
	switch commandName {
	case "start":
		HandleStart(bot, update)
	default:
		log.Printf("unknown command: %v\n", commandName)
	}
}

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message
	chatID := message.Chat.ID

	log.Printf("Получено сообщение: %v, от пользователя: %v", message.Text, message.From.UserName)

	chatState := chatstate.GetInstance()
	chatState.InitState(chatID, "start")

	keyboard := CreateKeyboard([]string{"Поиск предмета", "предметы на отслеж."}, 2)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я бот для поиска и отслеживания цен предметов в игре Black Desert Online")

	msg.ReplyMarkup = keyboard
	bot.Request(msg)
}

func StateRouter(bot *tgbotapi.BotAPI, update tgbotapi.Update, state string, indexMC int, indexSC int) {
	var message string
	// chatID := update.CallbackQuery.Message.Chat.ID
	var buttons []string
	var keyboard tgbotapi.InlineKeyboardMarkup
	// items := FillItems()

	var switchRout bool
	var itemrout bool
	
	switch state {
	case "start":
		keyboard = CreateKeyboard(startkeyboard, 2)
		message = "Главное меню"

	case "search":
		keyboard = CreateKeyboard([]string{"Осн. оружие", "Броня", "Аксессуары", "Назад"}, 2)
		message = "Выберите категорию предмета"

	case "weapon":
		keyboard = CreateKeyboard([]string{"Меч", "Назад"}, 2)
		message = "Вы выбрали категорию: Оружие"

	case "MainCRouting":
		switch indexMC{
		case 1:
			buttons, message = itemrouting.MainCRouting(1)
			keyboard = CreateKeyboard(buttons, 2)
		case 2:
			buttons, message = itemrouting.MainCRouting(2)
			keyboard = CreateKeyboard(buttons, 2)
		case 3:
			buttons, message = itemrouting.MainCRouting(3)
			keyboard = CreateKeyboard(buttons, 2)
		}
	case "SubCRouting":
		
		itemrout = true
		return
	case "switchRouting":

		switchRout = true
	}

	if itemrout != true && switchRout != true {
		EditMessage(update, bot, message, keyboard)
	}
}



func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	data := strings.ToLower(update.CallbackQuery.Data)
	chatState := chatstate.GetInstance() // Get the singleton instance

	var nextState string
	isPrev := false
	var indexMC int = 0
	var indexSC int = 0

	switch data {
	//--------------------------------StartRotate--------------------------------
	case "поиск предмета_callback":
		nextState = "search"
	case "предметы на отслеж._callback":


	//--------------------------------MainCategories--------------------------------
	case "осн. оружие_callback":
		nextState = "MainCRouting"
		indexMC = 1

	case "броня_callback":
		nextState = "MainCRouting"
		indexMC = 2
	case "бижутерия_callback":
		nextState = "MainCRouting"
		indexMC = 3

	//--------------------------------SubCatogories--------------------------------
	case "меч_callback":
		nextState = "itemsRouting"
		indexSC = 1
	case "лук_callback":
		nextState = "itemsRouting"
		indexSC = 1
	case "Шлем":
		nextState = "itemsRouting"
		indexSC = 2
	case "Доспехи":
		nextState = "itemsRouting"
		indexSC = 2
	case "Кольцо":
		nextState = "itemsRouting"
		indexSC = 3
	case "Ожерелье":
		nextState = "itemsRouting"
		indexSC = 3
	//--------------------------------switchStates--------------------------------
	case "предыдущий_callback":
		nextState = "switchRouting"

	case "следующий_callback":
		nextState = "switchRouting"

	case "назад_callback":
		nextState = chatState.PopState(chatID)
		isPrev = true
	default:
		log.Printf("unknown callback data: %v\n", data)
	}

	if len(nextState) != 0 {
		if !isPrev {
			chatState.PushState(chatID, nextState)
		}
		StateRouter(bot, update, nextState, indexMC, indexSC)
	}
}


func ItemRouting(bot *tgbotapi.BotAPI, update tgbotapi.Update, mainC int, subC int) {
	chatID := update.CallbackQuery.Message.Chat.ID
	callbackdata := update.CallbackQuery.Data 

	_, exists := categoryCashe[chatID]
	if !exists {
		categoryCashe[chatID] = struct{
			mainC int
			subC int
		}{
			mainC: mainC,
			subC: subC,
		}
	}

	category := categoryCashe[chatID]

	items, exists := itemCache[chatID]
	if !exists || len(items) == 0 {
		log.Println("items cache is empty")
		newItems, err := FillItems(chatID, category.mainC, category.subC )
		if err != nil {
			fmt.Println("Error to fillItems")
			return 
		}
		itemCache[chatID] = newItems
		items = newItems
	}

	var newIndex int

	if update.CallbackQuery != nil {
		switch callbackdata {
		case "Предыдущий_callback":
			newIndex = defineItemIndex(chatID, "prev")
		case "Следующий_callback":
			newIndex = defineItemIndex(chatID, "next")
		default:
			ItemOP(bot, update, items, newIndex)
			return
		}
		
		ItemOP(bot, update, items, newIndex)
		
	}
}

func defineItemIndex(chatID int64, direction string) int {
	curIndex, exists := itemIndexMap[chatID]
	if !exists {
		curIndex = 0
	}

	switch direction {
	case "prev":
		if curIndex > 0 {
			curIndex--
		}
	case "next":
		if curIndex < 0 {
			curIndex++
		}	
	} 

	itemIndexMap[chatID] = curIndex
	return curIndex
}

func ItemOP(bot *tgbotapi.BotAPI, update tgbotapi.Update, items []Item, curIndex int) {
	defer utils.TimeIt(time.Now(), "ItemOP") // Start timer

	if curIndex < 0 || curIndex >= len(items) {
		log.Println("Invalid item index")
		return
	}

	currentItem := items[curIndex]

	bodyMsg := fmt.Sprintf("Id предмета: %v, \nНазвание предмета: %v, \nЦена предмета: %v", currentItem.ID, currentItem.Name, currentItem.Price)

	keyboard := CreateKeyboard(switchKeyboard, 3)

	EditMessage(update, bot, bodyMsg, keyboard)
}

func FillItems(chatID int64, mainC int, subC int) ([]Item, error) {
	defer utils.TimeIt(time.Now(), "FillItems") // Start timer

	
	bdoItems, err := ba.GetWorldMarketList(mainC, subC)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(bdoItems))
	var wg sync.WaitGroup
	errChan := make(chan error, len(bdoItems))

	for i, bdoItem := range bdoItems {
		wg.Add(1)

		go func(i int, bdoItem ba.Item) {
			defer wg.Done()

			latestPrice, err := ba.GetLatestPrice(bdoItem.ID, 0)
			if err != nil {
				log.Printf("Failed to get price for item %d: %v", bdoItem.ID, err)
				errChan <- err
				return
			}

			items[i] = Item{
				ID:    bdoItem.ID,
				Name:  bdoItem.Name,
				Price: latestPrice,
			}
		}(i, bdoItem)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, fmt.Errorf("some items failed to load")
	}

	if len(items) < 1 {
		log.Println("error fill items")
		return nil, fmt.Errorf("error fill items")
	}

	itemCache[chatID] = items
	itemIndexMap[chatID] = 0
	return items, nil

}

var startkeyboard = []string{
	"Поиск предмета", "Отслеж. Товары",
}

var switchKeyboard = []string{
	"Предыдущий", "отслеживание", "Следующий", "Назад",
}

func CreateKeyboard(buttons []string, buttonsPerRow int) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton
	for i, btnText := range buttons {
		currentRow = append(currentRow, tgbotapi.NewInlineKeyboardButtonData(btnText, btnText+"_callback"))
		if (i+1)%buttonsPerRow == 0 {
			rows = append(rows, currentRow)
			currentRow = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}

func EditMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, message string, keyboard tgbotapi.InlineKeyboardMarkup) {
	var err error
	editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message)
	editMarkup := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		keyboard)
	_, err = bot.Request(editMsg)
	if err != nil {
		log.Printf("Ошибка при изменении сообщения: %v", err)
		return
	}
	_, err = bot.Request(editMarkup)
	if err != nil {
		log.Printf("Ошибка при изменении сообщения: %v", err)
	}

}