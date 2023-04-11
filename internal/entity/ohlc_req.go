package entity

import (
	"strconv"
	"strings"
)

type OHLCTransactionReq struct {
	Type             string `json:"type"`
	StockCode        string `json:"stock_code"`
	OrderBook        string `json:"order_book"`
	OrderNumber      string `json:"order_number"`
	Quantity         string `json:"quantity"`
	Price            string `json:"price"`
	ExecutedQuantity string `json:"executed_quantity"`
	ExecutionPrice   string `json:"execution_price"`
}

func (r OHLCTransactionReq) Translate() OHLCTransaction {
	var (
		price    = 0
		quantity = 0
	)

	if strings.TrimSpace(r.Quantity) != "" {
		quantity, _ = strconv.Atoi(r.Quantity)
	} else if strings.TrimSpace(r.ExecutedQuantity) != "" {
		quantity, _ = strconv.Atoi(r.ExecutedQuantity)
	}

	if strings.TrimSpace(r.Price) != "" {
		price, _ = strconv.Atoi(r.Price)
	} else if strings.TrimSpace(r.ExecutionPrice) != "" {
		price, _ = strconv.Atoi(r.ExecutionPrice)
	}

	return OHLCTransaction{
		Type:      r.Type,
		StockCode: r.StockCode,
		Quantity:  quantity,
		Price:     price,
	}
}
