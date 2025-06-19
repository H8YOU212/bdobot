package handlers

import (
	// "fmt"
	b "bdobot/bdoapi"
	"bdobot/db"
	"log"
	"strconv"
	"strings"

	// "bdobot/db"
	// "bdobot/tgbot/auth"
	"bdobot/tgbot/chatstate"
	itemrouting "bdobot/tgbot/itemRouting"
	isr "bdobot/tgbot/itemSpecRouting"

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

type MessageHandler func(*tgbotapi.Message)

var (
	itemIndexMap  = make(map[int64]int)
	itemCache     = make(map[int64][]Item)
	categoryCashe = make(map[int64]struct {
		mainC int
		subC  int
	})
)
var sid = new(int)
var targetprice int
var item = new(b.Item)
var indexMC = new(int)
var indexSC = new(int)
var curIndex = new(int)
var curItem = new(db.ItemSpec)

func HandleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {

	var err error
	var chatID int64
	var commandName string
	
	chatState := chatstate.GetInstance()

	if update.Message != nil { // Проверяем, что это текстовое сообщение
		chatID = update.Message.Chat.ID
		commandName = update.Message.Command()
	} else if update.CallbackQuery != nil { // Проверяем, что это callback-запрос
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		log.Println("Неизвестный тип обновления")
		return
	}
	lastState := chatState.GetLastState(chatID)

	baseUser := db.User{
		ID:          chatID,
		Name:        update.Message.From.UserName,
		ItemsOnSpec: []db.ItemSpec{},
	}

	switch commandName {
	case "start":
		HandleStart(bot, update)
		exists, err := db.UserExists(chatID)
		if err != nil {
			log.Printf("Ошибка проверки пользователя: %v", err)
			return
		}
		if !exists {
			db.InsertUser(baseUser)
		}
	default:
		log.Printf("Неизвестная команда: %v\n", commandName)
	}


	if lastState == "AddSpecItem" {
		messageText := update.Message.Text
		targetprice, err = strconv.Atoi(messageText)
		if err != nil {
			log.Println("Ошибка преобразования str to int")
			EditMessageWoutMarkup(update, bot, "Вы ввели не число")
			return
		}
		chatState.PushState(chatID, "setTargetPrice")
		nextState := "setTargetPrice"
		HandlePrice(bot, update, nextState, targetprice, chatID)
		chatState.PopState(chatID)
		// StateRouter(bot, update, "search", indexMC, indexSC, curIndex, chatID)
	}

	if lastState == "setSid" {
		messageText := update.Message.Text
		settedSid, err := strconv.Atoi(messageText)
		if err != nil {
			log.Println("Ошибка преобразования str to int")
			EditMessageWoutMarkup(update, bot, "Вы ввели не число")
			return
		} else if settedSid >= 20 {
			EditMessageWoutMarkup(update, bot, "Уровень предмета не может быть выше 20")
			return
		}
		*sid = settedSid
		chatState.PopState(chatID)
		// var message string
		// buttons, message, (*item) = itemrouting.SubCRouting(*indexMC, *indexSC, curIndex, *sid)
		// keyboard = CreateKeyboard(buttons, 3)
		// StateRouter(bot, update, "search", indexMC, indexSC, curIndex, chatID)
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

func HandlePrice(bot *tgbotapi.BotAPI, update tgbotapi.Update, state string, targetPrice int, chatID int64) {

	isr.AddNewItem(chatID, targetPrice, *item)

}

func StateRouter(bot *tgbotapi.BotAPI, update tgbotapi.Update, state string, indexMC *int, indexSC *int, curIndex *int, chatID int64) {
	var message string
	// chatID = update.CallbackQuery.Message.Chat.ID
	var buttons []string
	var keyboard tgbotapi.InlineKeyboardMarkup
	// items := new([]ba.Item)

	chatState := chatstate.GetInstance()
	lastState := chatState.GetLastState(chatID)

	// curIndex := new(int)
	var switchRout bool
	var itemrout bool

	switch state {
	//afstart
	case "start":
		keyboard = CreateKeyboard(startkeyboard, 2)
		message = "Привет! Я бот для поиска и отслеживания цен предметов в игре Black Desert Online"

	case "search":
		keyboard = CreateKeyboard([]string{"Осн. оружие", "Броня", "Аксессуары", "Назад"}, 2)
		message = "Выберите категорию предмета"

	case "specItems":
		keyboard = CreateKeyboard([]string{"Предыдущий", "Удалить", "Следующий", "Назад"}, 3)
		message, *curItem = isr.GetSpcItms(chatID, curIndex)

	//ItemRouting
	case "MainCRouting":
		buttons, message = itemrouting.MainCRouting(*indexMC)
		keyboard = CreateKeyboard(buttons, 2)

	case "SubCRouting":
		*curIndex = 0
		buttons, message, (*item) = itemrouting.SubCRouting(*indexMC, *indexSC, curIndex, *sid)
		keyboard = CreateKeyboard(buttons, 3)

	// SwitchRouting
	case "switchRouting":
		if lastState == "SubCRouting" {
			buttons, message, (*item) = itemrouting.SubCRouting(*indexMC, *indexSC, curIndex, *sid)
			// switchRout = true
			keyboard = CreateKeyboard(buttons, 3)
		} else if lastState == "specItems" {
			keyboard = CreateKeyboard([]string{"Предыдущий", "Удалить", "Следующий", "Назад"}, 3)
			message, *curItem = isr.GetSpcItms(chatID, curIndex)
		}

	case "setSid":
		itemrout = true
		EditMessageWoutMarkup(update, bot, "Введите желаемый уровень предмета не выше 20-ого \n(Может повлиять на отображение цены предмета)")

	case "AddSpecItem":
		itemrout = true
		// message = "Введите грейд предмета"
		// EditMessageWoutMarkup(update, bot, message)

		EditMessageWoutMarkup(update, bot, "Введите желаемую цену для уведомления")

	case "deleteSpecItem":
		itemrout = true
		isr.DeleteSpecItem(chatID, curItem)
	}

	if itemrout != true && switchRout != true {
		EditMessage(update, bot, message, keyboard)
		return
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

	//--------------------------------MainCategories-----------------------------
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

	//--------------------------------SubCatogories------------------------------
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

	//--------------------------------settingItems--------------------------------
	case "установить sid_callback":
		nextState = "setSid"

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
	case "отслеживание_callback":
		nextState = "AddSpecItem"

	case "удалить_callback":
		nextState = "deleteSpecItem"
		isPrev = true

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
		StateRouter(bot, update, nextState, indexMC, indexSC, curIndex, chatID)
	}
}

var startkeyboard = []string{
	"Поиск предмета", "предметы на отслеж.",
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

    if update.CallbackQuery == nil  {
        log.Printf("EditMessage: CallbackQuery = nil")
        return
    }

	if update.CallbackQuery.Message == nil {
        log.Printf("EditMessage: Message равен nil")
        return
    }

	if update.CallbackQuery.Message.Chat == nil {
        log.Printf("EditMessage: Chat равен nil")
        return
    }

    editMsg := tgbotapi.NewEditMessageText(
        update.CallbackQuery.Message.Chat.ID,
        update.CallbackQuery.Message.MessageID,
        message,
    )
    editMarkup := tgbotapi.NewEditMessageReplyMarkup(
        update.CallbackQuery.Message.Chat.ID,
        update.CallbackQuery.Message.MessageID,
        keyboard,
    )
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

func EditMessageWoutMarkup(update tgbotapi.Update, bot *tgbotapi.BotAPI, message string) {
	var err error
	editMsg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		message)
	_, err = bot.Request(editMsg)
	if err != nil {
		log.Printf("Ошибка при изменении сообщения: %v", err)
		return
	}
}
