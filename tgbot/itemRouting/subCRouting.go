package itemrouting

import b "bdobot/bdoapi"



func SubCRouting(mainC int, subC int, curIndex *int) ([]string, string, b.Item) {
	var keyboard []string
	var message string
	// var item b.Item
	message, item := DefineItem(curIndex, mainC, subC)

	keyboard = []string{
		"Предыдущий", "отслеживание", "Следующий", "Назад",
	}
	return keyboard, message, item

}
