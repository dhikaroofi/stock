package entity

// NewRecords is struct for contain new record
type NewRecords struct {
	StockCode string
	Price     int64
	Quantity  int64
}

// IndexMember is struct for contain index member of stock
type IndexMember struct {
	StockCode string
	IndexCode string
}

// Summary is struct for contain summary data from OHLC
type Summary struct {
	StockCode string   `json:"stock_code"`
	IndexCode []string `json:"index_code"`
	Open      int64    `json:"open"`
	High      int64    `json:"high"`
	Low       int64    `json:"low"`
	Close     int64    `json:"close"`
	Prev      int64    `json:"prev"`
}
