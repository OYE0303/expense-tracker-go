package transaction

import "time"

type icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type mainCateg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Icon icon   `json:"icon"`
}

type subCateg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type transaction struct {
	ID        int64     `json:"id"`
	MainCateg mainCateg `json:"main_category"`
	SubCateg  subCateg  `json:"sub_category"`
	Price     float64   `json:"price"`
	Note      string    `json:"note"`
	Date      time.Time `json:"date"`
}

type getTransactionResp struct {
	Transactions []transaction `json:"transactions"`
}