package entity

type ListOHLCTransactions struct {
	Info string
	List []OHLCTransaction `json:"list"`
}

// OHLCTransaction is used for parsing data from data streamer to the core use case
type OHLCTransaction struct {
	Type      string `json:"type"`
	StockCode string `json:"stock_code"`
	OrderBook string `json:"order_book"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
}

type ResultSummary struct {
	Result map[string]OHLCSummary `json:"result"`
}

// OHLCSummary is used for show the result form OHLC calculation
type OHLCSummary struct {
	Date          string  `json:"date"`
	StockCode     string  `json:"stock_code"`
	Average       float64 `json:"average"`
	AverageRound  float64 `json:"average_round"`
	Volume        int64   `json:"volume"`
	Value         int64   `json:"value"`
	PreviousPrice int     `json:"previous_price"`
	HighestPrice  int     `json:"highest_price"`
	LowestPrice   int     `json:"lowest_price"`
	ClosePrice    int     `json:"close_price"`
	TotalTrans    int     `json:"total_trans"`
}
