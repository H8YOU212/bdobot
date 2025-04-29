package itemspecrouting

import (
	b "bdobot/bdoapi"
	"bdobot/db"
)

func AddNewItem(id int64, targetprice int, item b.Item) {
	itemSpec := db.ItemSpec{
		ID: item.ID,
		Name: item.Name,
		Price: item.Price,
		ItemStartPrice: item.Price,
		ItemTargetPrice: targetprice,
	}
	
	db.AddItemToSpecItems(id, itemSpec)

}
