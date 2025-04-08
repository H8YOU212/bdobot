package itemrouting

import(
	// b "bdobot/bdoapi"
)

func SubCRouting(mainC int, subC int, curIndex *int) ([]string, string) {
	var keyboard []string
	var message string
	message = DefineItem(curIndex, mainC, subC)

	keyboard = []string{
		"Предыдущий", "отслеживание", "Следующий", "Назад",
	}
	return keyboard, message
	
}