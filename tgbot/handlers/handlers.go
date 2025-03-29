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

func StateRouter(bot *tgbotapi.BotAPI, update tgbotapi.Update, state string) {
	var message string
	var keyboard tgbotapi.InlineKeyboardMarkup

	switch state {
	case "start":
		keyboard = CreateKeyboard(startkeyboard, 2)
		message = "Главное меню"

	case "search":
		keyboard = CreateKeyboard([]string{"Осн. оружие", "Броня", "Аксессуары", "Назад"}, 2)
		message = "Выберите категорию предмета"

	case "armour":
		
	case "weapon":
		keyboard = CreateKeyboard([]string{"Меч", "Назад"}, 2)
		message = "Вы выбрали категорию: Оружие"
	case "sword":
		// 	keyboard = CreateKeyboard([]string{"Меч Элси", "Назад"}, 2)
		// 	message = "Выберите предмет:"

	}

	EditMessage(update, bot, message, keyboard)
}



func HandleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID
	data := strings.ToLower(update.CallbackQuery.Data)
	chatState := chatstate.GetInstance() // Get the singleton instance
	

	var nextState string
	isPrev := false

	switch data {
	case "поиск предмета_callback":
		nextState = "search"
	case "осн. оружие_callback":
		nextState = "weapon"
	case "меч_callback":
		nextState = "sword"
		ItemRouting( bot, update, 1, 1) // itemRouting заменить функционал на ItemOP, и переместить в handleMessage, там же сдеать обработку пред, след, и обрабатывать состояния в виде map например: chatstate[subC]int, sword, текущее состояние подкатегории, int индекс предмета, который будет соответственно меняться при обработке пред, след
		return
	case "броня_callback":
		nextState = "armour"
	case "назад_callback":
		nextState = chatState.PopState(chatID)
		isPrev = true
	case "предыдущий_callback":
		
	case "следующий_callback":
		

	default:
		log.Printf("unknown callback data: %v\n", data)
	}

	if len(nextState) != 0 {
		if !isPrev {
			chatState.PushState(chatID, nextState)
		}
		StateRouter(bot, update, nextState)
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

// var mainCategories = []string{
// 	""
// }

var subCategories = []string{
	"Меч", "лук",
}

var switchKeyboard = []string{
	"Предыдущий", "отслеживание", "Следующий", "Назад",
}

func fillKeyboard() ([]string, int) {
	panic("unimplemented")
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

func CreateMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, message string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msg.ReplyMarkup = keyboard
	bot.Request(msg)
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

func (b *Bot) AddHandler(command string, handler MessageHandler) {
	b.handlers[command] = handler
}
