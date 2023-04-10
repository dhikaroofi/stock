package ohlc

import (
	"context"
	cahce "github.com/dhikaroofi/stock.git/internal/adapters/driven/cache"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"math"
	"strings"
)

type Task interface {
	Calculate(ctx context.Context, info string, transactions []entity.OHLCTransaction) error
	GetSummary(ctx context.Context, id string) error
}

type ohlcUseCase struct {
	cache cahce.Task
}

func New(cache cahce.Task) Task {
	return &ohlcUseCase{
		cache: cache,
	}
}

func (uc ohlcUseCase) Calculate(ctx context.Context, info string, transactions []entity.OHLCTransaction) error {
	log := logger.ExtractLogFromContext(ctx)

	date := info[0:10]

	listStock := make(map[string]int)
	for _, val := range transactions {
		if _, ok := listStock[val.StockCode]; !ok {
			listStock[val.StockCode] = 0
		}
		listStock[val.StockCode] = listStock[val.StockCode] + 1
	}

	listSummaryStock := make(map[string]entity.OHLCSummary, len(listStock))
	for _, val := range transactions {
		var summary entity.OHLCSummary
		if summaryStock, ok := listSummaryStock[val.StockCode]; ok {
			summary = summaryStock
		}

		if val.Quantity == 0 {
			summary.PreviousPrice = val.Price
			summary.HighestPrice = val.Price
			summary.LowestPrice = val.Price
		}

		if valType := strings.ToLower(strings.TrimSpace(val.Type)); valType == "e" || valType == "p" {
			if summary.TotalTrans < 1 {
				summary.HighestPrice = val.Price
				summary.LowestPrice = val.Price
			}

			if val.Price > summary.HighestPrice {
				summary.HighestPrice = val.Price
			}
			if val.Price > 0 && (val.Price < summary.LowestPrice) {
				summary.LowestPrice = val.Price
			}
			if summary.TotalTrans >= listStock[val.Type] {
				summary.ClosePrice = val.Price
			}

			summary.Volume = summary.Volume + int64(val.Quantity)
			summary.Value = summary.Value + int64(val.Quantity*val.Price)
			summary.TotalTrans++
		}

		if summary.TotalTrans >= listStock[val.Type] {
			if summary.Volume > 0 {
				summary.Average = float64(summary.Value) / float64(summary.Volume)
			}
		}

		listSummaryStock[val.StockCode] = summary
	}

	for index, _ := range listSummaryStock {
		summary, err := uc.cache.GetOHLCSummary(ctx, index, date)
		if err != nil {
			log.SetError(err).Info(ctx, "failed to get summary from redis")
			continue
		}

		if summary.PreviousPrice == 0 {
			summary.PreviousPrice = listSummaryStock[index].PreviousPrice
		}
		if summary.ClosePrice == 0 {
			summary.ClosePrice = listSummaryStock[index].ClosePrice
		}
		if summary.LowestPrice == 0 {
			summary.LowestPrice = listSummaryStock[index].LowestPrice
		}

		if listSummaryStock[index].HighestPrice > summary.HighestPrice {
			summary.HighestPrice = listSummaryStock[index].HighestPrice
		}

		if listSummaryStock[index].LowestPrice > 0 && (listSummaryStock[index].LowestPrice < summary.LowestPrice) {
			summary.LowestPrice = listSummaryStock[index].LowestPrice
		}

		summary.Value = summary.Value + listSummaryStock[index].Value
		summary.Volume = summary.Volume + listSummaryStock[index].Volume
		summary.TotalTrans = summary.TotalTrans + listSummaryStock[index].TotalTrans
		if summary.Value > 0 {
			summary.Average = float64(summary.Value) / float64(summary.Volume)
			summary.AverageRound = math.Round(summary.Average)
		}
		listSummaryStock[index] = summary
	}

	if err := uc.cache.SaveOHLCSummaryBatch(ctx, date, listSummaryStock); err != nil {
		log.SetError(err).Info(ctx, "can't save summary to redis")
		return err
	}

	return nil
}

func (uc ohlcUseCase) GetSummary(ctx context.Context, id string) error {
	return nil
}
