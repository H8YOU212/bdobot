package notification

import "fmt"


func fake_notify(user *User) (string, []string) {
	var items []Item

	items = []Item{
		{
			Name: "Меч Асвол",
			ID:   10007,
		},
	}

	msg := fmt.Sprintf("Предмет: Название: %v\nid: %v\nPrice: %v\n\nДостиг или превысил установленную сумму", items[0].Name,items[0].ID)
	return msg, []string{"Назад"}
}