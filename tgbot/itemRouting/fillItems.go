package itemrouting

import (
	b "bdobot/bdoapi"
	"fmt"
	"log"
	"sync"
)

func FillItems(chatID int64, mainC int, subC int) ([]b.Item, error) {
	bdoItems, err := b.GetWorldMarketList(mainC, subC)
	if err != nil {
		return nil, err
	}

	items := make([]b.Item, len(bdoItems))
	var wg sync.WaitGroup
	errChan := make(chan error, len(bdoItems))

	for i, bdoItem := range bdoItems {
		wg.Add(1)

		go func(i int, bdoItem b.Item) {
			defer wg.Done()

			latestPrice, err := b.GetLatestPrice(bdoItem.ID, 0)
			if err != nil {
				log.Printf("Failed to get price for item %d: %v", bdoItem.ID, err)
				errChan <- err
				return
			}

			items[i] = b.Item{
				ID:    bdoItem.ID,
				Name:  bdoItem.Name,
				Price: latestPrice,
			}
		}(i, bdoItem)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, fmt.Errorf("some items failed to load")
	}

	if len(items) < 1 {
		log.Println("error fill items")
		return nil, fmt.Errorf("error fill items")
	}

	// itemCache[chatID] = items TODO: realization itemCashe func 
	itemIndexMap[chatID] = 0
	return items, nil
}
