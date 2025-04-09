package itemrouting

func MainCRouting(mainC int) ([]string, string) {
	var keyboard []string
	var message string
	switch mainC {
	case 1:
		keyboard = []string{"меч", "лук", "Назад"}
		message = "Вы выбрали категорию: Оружие"

	case 5: 
		keyboard = []string{"Щит", "Кинжал", "Назад"}
		message = "Вы выбрали категорию: Доп. Оружие"		
		
	case 10:
		keyboard = []string{"Двуручный меч", "Коса", "Назад"}
		message = "Вы выбрали категориюю: Пробужденное оружие"

	case 15:
		keyboard = []string{"шлем", "доспехи", "Назад"}
		message = "Вы выбрали категорию: Броня"
		
	case 20:
		keyboard = []string{"Кольцо", "Ожерелье", "Назад"}
		message = "Вы выбрали категорию: Аксессуары"

	}
	return keyboard, message
	
}