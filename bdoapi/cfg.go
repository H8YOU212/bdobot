package bdoapi

type Item struct {
	ID    int    `json:"id"`
	Sid   int    `json:"sid"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	// BasePrice    int            `json:"basePrice"`
	// History      map[string]int `json:"history"`
	// MainCategory int            `json:"mainCategory"`
	// SubCategory  int            `json:"subCategory"`
	// PriceMin     int            `json:"priceMin"`
	// PriceMax     int            `json:"priceMax"`
}

type MarketPriceInfo struct {
	Name    string         `json:"name"`
	ID      int            `json:"id"`
	Sid     int            `json:"sid"`
	History map[string]int `json:"history"` // История цен в виде карты (timestamp -> цена)
}