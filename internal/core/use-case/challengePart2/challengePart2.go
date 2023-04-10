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

type challengePart2 struct{}

func (challengePart2) OHLCUpdater(ctx context.Context, ohlc entity.OHLCIndexMember) (map[string]entity.Summary, error) {
	var (
		result = map[string]entity.Summary{}
		record entity.Summary
	)

	for _, stock := range ohlc.Stock {
		result[stock] = entity.Summary{StockCode: stock}
	}

	for _, val := range ohlc.NewRecords {
		record = result[val.StockCode]
		switch {
		case val.Quantity == 0:
			record.Prev = val.Price
		case val.Quantity > 0 && record.Open == 0:
			record.Open = val.Price
		default:
			record.Close = val.Price
			if record.High < val.Price {
				record.High = val.Price
			}
			if record.Low > val.Price {
				record.Low = val.Price
			}
		}

		result[val.StockCode] = record
	}

	for _, val := range ohlc.IndexMember {
		record = result[val.StockCode]
		record.IndexCode = append(record.IndexCode, val.IndexCode)
		result[val.StockCode] = record
	}

	for _, val := range result {
		payload, _ := json.Marshal(val)
		fmt.Println(string(payload))
	}

	return result, nil
}
