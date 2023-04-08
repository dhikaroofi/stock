package entity

type OHLCTransaction struct {
	Type      string `json:"type"`
	StockCode string `json:"stock_code"`
	OrderBook int    `json:"order_book"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

type OHLCSummary struct {
	PreviousPrice int `json:"type"`
	OpenPrice     int `json:"stock_code"`
	HighestPrice  int `json:"order_book"`
	LowestPrice   int `json:"quantity"`
	ClosePrice    int `json:"price"`
	Volume        int `json:"volume"`
	Value         int `json:"value"`
	Average       int `json:"average"`
	AverageRound  int `json:"average_round"`
	TotalTrans    int `json:"total_trans"`
}
