package itemspecrouting

import (
	b "bdobot/bdoapi"
	"bdobot/db"
)

func AddNewItem(id int64, targetprice int, item b.Item) {
	var method int
	if item.Price > targetprice {
		method = 0
		itemSpec := db.ItemSpec{
			ID:              	item.ID,
			SID:				item.Sid,
			Name:            	item.Name,
			Price:           	item.Price,
			ItemStartPrice:  	item.Price,
			ItemTargetPrice: 	targetprice,
			Method:          	method,
		}
		db.AddItemToSpecItems(id, itemSpec)
	} else {
		method = 1
		itemSpec := db.ItemSpec{
			ID:              item.ID,
			SID:				item.Sid,
			Name:            item.Name,
			Price:           item.Price,
			ItemStartPrice:  item.Price,
			ItemTargetPrice: targetprice,
			Method:          method,
		}
		db.AddItemToSpecItems(id, itemSpec)
	}

}
