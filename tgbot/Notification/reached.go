package notification

import "bdobot/db"

type User db.User

type Item db.ItemSpec

func Notify(user *User) {
	items := user.ItemsOnSpec
	for _, i := range items {
		
	}
}
