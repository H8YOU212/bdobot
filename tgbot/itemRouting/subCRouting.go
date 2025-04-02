package itemrouting

import(
	b "bdobot/bdoapi"
)

func SubCRouting(subC int, items []b.Item) ([]string, string) {
	var keyboard []string
	var message string
	switch subC {
	case 1:
		keyboard = []string{"меч", "лук", "Назад"}
		message = "Вы выбрали категорию: Оружие"		
		
	case 2:
		keyboard = []string{"шлем", "доспехи", "Назад"}
		message = "Вы выбрали категорию: Броня"
		
	case 3:
		keyboard = []string{"Кольцо", "Ожерелье", "Назад"}
		message = "Вы выбрали категорию: Аксессуары"

	}
	return keyboard, message
	
}