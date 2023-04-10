package challengepart2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dhikaroofi/stock.git/internal/entity"
)

// Task is an interface that contain several task that has in challenge part 2
type Task interface {
	OHLCUpdater(ctx context.Context, ohlc entity.OHLCIndexMember) (map[string]entity.Summary, error)
}

// New is constructor for challenge part 2
func New() Task {
	return &challengePart2{}
}

type challengePart2 struct {
}

func (challengePart2) OHLCUpdater(ctx context.Context, ohlc entity.OHLCIndexMember) (map[string]entity.Summary, error) {
	var result = map[string]entity.Summary{}

	for _, stock := range ohlc.Stock {
		var stockSummary entity.Summary
		if val, ok := result[stock]; ok {
			stockSummary = val
		} else {
			stockSummary.StockCode = stock
		}

		updateRecord(&stockSummary, ohlc.NewRecords)

		for _, val := range ohlc.IndexMember {
			if val.StockCode == stock {
				stockSummary.IndexCode = append(stockSummary.IndexCode, val.IndexCode)
			}
		}

		result[stock] = stockSummary
		payload, _ := json.Marshal(stockSummary)
		fmt.Println(string(payload))
	}
	return result, nil
}

// This loop updates the stock data with the new records.
// The updated data includes quantity, previous price, opening price, high price, and low price.
func updateRecord(stockSummary *entity.Summary, newRecords []entity.NewRecords) {
	for _, record := range newRecords {
		if record.StockCode == stockSummary.StockCode {
			switch {
			case record.Quantity == 0:
				stockSummary.Prev = record.Price
			case record.Quantity > 0 && stockSummary.Open == 0:
				stockSummary.Open = record.Price
			default:
				stockSummary.Close = record.Price
				if stockSummary.High < record.Price {
					stockSummary.High = record.Price
				}
				if stockSummary.Low > record.Price {
					stockSummary.Low = record.Price
				}
			}
		}
	}
}
