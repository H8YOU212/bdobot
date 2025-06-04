package itemspecrouting

import (
	"bdobot/db"
	"fmt"
)

func GetSpcItms(id int64, curIndex *int) (string, db.ItemSpec) {
	var message string
	user := GetUser(id)
	if user == nil {
		fmt.Println("Undefine user by id(check itemSpecRouting/getSpecItems)")
		return "", db.ItemSpec{}
	}

	specitems := user.ItemsOnSpec

	if *curIndex >= 0 && *curIndex < len(specitems) {
		item := specitems[*curIndex]
		message = fmt.Sprintf(fmt.Sprintf("Id предмета: %v, \nНазвание предмета: %v, \nЦена предмета: %v, \nНачальная Цена: %v, \nПлановая цена %v", item.ID, item.Name, item.Price, item.ItemStartPrice, item.ItemTargetPrice))
		return message, item
	}
	return "", db.ItemSpec{}
}
