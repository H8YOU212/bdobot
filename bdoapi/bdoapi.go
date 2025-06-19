package bdoapi

import (
	logfdb "bdobot/logger"
	"bdobot/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	// "strings"
)

const baseUrl = "https://api.arsha.io"



func GetWorldMarketList(mainCategory int, subCategory int) ([]Item, error) {
	defer utils.TimeIt(time.Now(), "GetWorldMarketList")
	url := baseUrl + fmt.Sprintf("/v2/ru/GetWorldMarketList?mainCategory=%d&subCategory=%d&lang=ru", mainCategory, subCategory)
	// payload := strings.NewReader(fmt.Sprintf("{\n\t\"mainCategory\": %v,\n\t\"subCategory\": %v,\n\t\"lang\": \"ru\"}", mainCategory, subCategory))
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "BlackDesert")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)
	fmt.Println(string(body))

	var items []Item

	err = json.Unmarshal(body, &items)
	if err != nil {
		fmt.Println("parse JSON data error :", err)
		return nil, err
	}
	// fmt.Println(items)
	return items, nil

}

func GetMarketPriceInfo(id int, sid int) (map[string]int, error) {
	url := baseUrl + fmt.Sprintf("/v2/ru/GetMarketPriceInfo?id=%d&sid=%d&lang=ru", id, sid) // ?id=12237&sid=0
	method := "GET"
	/*payload := strings.NewReader(fmt.Sprintf(`{
		"id": %d,
		"sid": %d
	}`, id, sid))*/

	req, _ := http.NewRequest(method, url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "BlackDesert")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		e := fmt.Errorf("ошибка при чтении тела ответа: %v", err)
		return nil, e
	}

	var MarketPriceInfo MarketPriceInfo
	err = json.Unmarshal(body, &MarketPriceInfo)
	if err != nil {
		fmt.Println("Error parse data")
		return nil, err
	}

	logfdb.Logapi(url, string(body))

	return MarketPriceInfo.History, nil

}

func GetWorldMarketHotList() ([]Item, error) {
	url := "https://api.arsha.io/v2/ru/GetWorldMarketHotList?lang=ru"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var items []Item
	err = json.Unmarshal(body, &items)
	if err != nil {
		return nil, err
	}

	return items, nil
}
