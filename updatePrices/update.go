package updateprices

import (
	"bdobot/bdoapi"
	"bdobot/db"
	"fmt"
	"log"
	"time"
)

type Item db.ItemSpec

type User db.User

func StartUpdater() ( error ){
	ticker := time.NewTicker(10 * time.Minute) // интервал обновления
	defer ticker.Stop()

	for {
		log.Println("Запуск обновления цен для всех пользователей...")
		users, err := db.GetAllUsers()
		if err != nil {
			log.Printf("Ошибка получения пользователей: %v", err)
			continue
		}

		for _, user := range users {
			// Здесь вызывайте функцию обновления цен для каждого пользователя
			changePrice(&user)
		}

		<-ticker.C // ждем следующего тика
	}
	
}

func changePrice(user *db.User) {
	userItems := user.ItemsOnSpec

	for _, i := range userItems {
		getid := i.ID
		getSid := i.SID
		up := i.Price
		curPrice, err := bdoapi.GetLatestPrice(getid, getSid)
		if err != nil {
			fmt.Println(err)
		}

		if up != curPrice {
			i.Price = curPrice
		}
	}

}
