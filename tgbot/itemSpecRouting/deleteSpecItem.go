package itemspecrouting

import "bdobot/db"

func DeleteSpecItem(id int64, itemSpec *db.ItemSpec) {
	db.Delete(id, *itemSpec)
}
