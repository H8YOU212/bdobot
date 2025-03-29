package bdoapi

import "fmt"

// func GetLastUpdate(item *Item) time.Time {

// }

func GetLatestPrice(id int, sid int) (int, error) {
	itemhistory, err := GetMarketPriceInfo(id, sid)
	if err != nil {
		fmt.Println("Error pull up history item data")
		return 0, err
	}

	if len(itemhistory) == 0 {
		return 0, fmt.Errorf("история цен пуста")
	}

	// Получаем последний ключ (максимальный timestamp)
	var latestTimestamp string
	for timestamp := range itemhistory {
		if latestTimestamp == "" || timestamp > latestTimestamp {
			latestTimestamp = timestamp
		}
	}

	// Возвращаем цену по последнему timestamp
	return itemhistory[latestTimestamp], nil
}
