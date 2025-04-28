package handlers

import (
	"log"
	"strings"
	
	"bdobot/tgbot/chatstate"
	"bdobot/tgbot/itemRouting"
	"bdobot/tgbot/auth"

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

type user auth.User

type MessageHandler func(*tgbotapi.Message)

var ( 
	itemIndexMap = make(map[int64]int)
	itemCache = make(map[int64][]Item)
	categoryCashe = make(map[int64]struct{
		mainC int
		subC int
	})
)


var indexMC = new(int)
var indexSC = new(int)
var curIndex = new(int)

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	commandName := update.Message.Command()
	switch commandName {
	case "start":
		HandleStart(bot, update)
	default:
		log.Printf("unknown command: %v\n", commandName)
	}

	if chatstate.
}

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message
	chatID := message.Chat.ID

	log.Printf("Получено сообщение: %v, от пользователя: %v", message.Text, message.From.UserName)

	chatState := chatstate.GetInstance()
	chatState.InitState(chatID, "start")

	keyboard := CreateKeyboard([]string{"Поиск предмета", "предметы на отслеж.", "Авторизация"}, 2)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я бот для поиска и отслеживания цен предметов в игре Black Desert Online")

	msg.ReplyMarkup = keyboard
	bot.Request(msg)
}

func StateRouter(bot *tgbotapi.BotAPI, update tgbotapi.Update, state string, indexMC *int, indexSC *int, curIndex *int) {
	var message string
	// chatID := update.CallbackQuery.Message.Chat.ID
	var buttons []string
	var keyboard tgbotapi.InlineKeyboardMarkup
	// items := new([]ba.Item)
	


	// curIndex := new(int)
	var switchRout bool
	var itemrout bool
	
	switch state {
	case "start":
		keyboard = CreateKeyboard(startkeyboard, 2)
		message = "Главное меню"

	case "search":
		keyboard = CreateKeyboard([]string{"Осн. оружие", "Броня", "Аксессуары", "Назад"}, 2)
		message = "Выберите категорию предмета"

	case "auth":
		
	case "specItems":
		
		
	
	case "MainCRouting":
		buttons, message = itemrouting.MainCRouting(*indexMC)
		keyboard = CreateKeyboard(buttons, 2)
	case "SubCRouting":
		buttons, message = itemrouting.SubCRouting(*indexMC,*indexSC, curIndex)
		keyboard = CreateKeyboard(buttons, 3)
	case "switchRouting":
		buttons, message =itemrouting.SubCRouting(*indexMC, *indexSC, curIndex)
		// switchRout = true
		keyboard = CreateKeyboard(buttons, 3)
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
	

	switch data {
	//--------------------------------StartRotate--------------------------------
	case "поиск предмета_callback":
		nextState = "search"
	case "предметы на отслеж._callback":
		nextState = "specItems"
	case "авторизация_callback":
		nextState = "auth"


	//--------------------------------MainCategories--------------------------------
	case "осн. оружие_callback":
		nextState = "MainCRouting"
		*indexMC = 1
	case "доп.оружие_callback":
		nextState = "MainCRouting"
		*indexMC = 5
	case "Пробужд. Оружие":
		nextState = "MainCRouting"
		*indexMC = 10
	case "броня_callback":
		nextState = "MainCRouting"
		*indexMC = 15
	case "аксессуары_callback":
		nextState = "MainCRouting"
		*indexMC = 20

	//--------------------------------SubCatogories--------------------------------
	case "меч_callback":
		nextState = "SubCRouting"
		*indexSC = 1
	case "лук_callback":
		nextState = "SubCRouting"
		*indexSC = 2
	case "щит_callback":
		nextState = "SubCRouting"
		*indexSC = 1
	case "кинжал_callback":
		nextState = "SubCRouting"
		*indexSC = 2
	case "двуручный меч_callback":
		nextState = "SubCRouting"
		*indexSC = 1
	case "коса_callback":
		nextState = "SubCRouting"
		*indexSC = 2
	case "шлем_callback":
		nextState = "SubCRouting"
		*indexSC = 1
	case "доспехи_callback":
		nextState = "SubCRouting"
		*indexSC = 2
	case "кольцо_callback":
		nextState = "SubCRouting"
		*indexSC = 1
	case "ожерелье_callback":
		nextState = "SubCRouting"
		*indexSC = 2
	//--------------------------------switchStates--------------------------------
	case "предыдущий_callback":
		nextState = "switchRouting"
		isPrev = true
		// if *curIndex >= 0{
			*curIndex--
		// }
	case "следующий_callback":
		nextState = "switchRouting"
		isPrev = true
		// if *curIndex <= 0{
			*curIndex++
		// }
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
		StateRouter(bot, update, nextState, indexMC, indexSC, curIndex)
	}
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