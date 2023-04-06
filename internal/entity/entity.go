package entity

type OHLCTransactions struct {
	Type      string `json:"type"`
	StockCode string `json:"stock_code"`
	OrderBook int    `json:"order_book"`
	Price     int    `json:"price"`
}
