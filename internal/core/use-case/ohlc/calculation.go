package ohlc

import (
	"context"
	"math"
	"strings"

	cahce "github.com/dhikaroofi/stock.git/internal/adapters/driven/cache"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/logger"
)

type Task interface {
	Calculate(ctx context.Context, transactions entity.ListOHLCTransactions) (entity.ResultSummary, error)
	GetSummary(ctx context.Context, date, stockCode string) (entity.OHLCSummary, error)
}

type ohlcUseCase struct {
	cache cahce.Task
}

func New(cache cahce.Task) Task {
	return &ohlcUseCase{
		cache: cache,
	}
}

func (uc ohlcUseCase) GetSummary(ctx context.Context, date, stockCode string) (entity.OHLCSummary, error) {
	summary, err := uc.cache.GetOHLCSummary(ctx, date, stockCode)
	if err != nil {
		return entity.OHLCSummary{}, err
	}
	return summary, nil
}

func (uc ohlcUseCase) Calculate(ctx context.Context,
	transactions entity.ListOHLCTransactions) (entity.ResultSummary, error) {

	log := logger.ExtractLogFromContext(ctx)
	date := transactions.Info[0:10]

	listSummaryStock := generateListSummaryStock(transactions.List).Result
	for index := range listSummaryStock {
		summary, err := uc.cache.GetOHLCSummary(ctx, index, date)
		if err != nil {
			log.SetError(err).Info(ctx, "failed to get summary from redis")
			continue
		}

		updateSummary(&summary, listSummaryStock[index])
		listSummaryStock[index] = summary
	}

	if err := uc.cache.SaveOHLCSummaryBatch(ctx, date, listSummaryStock); err != nil {
		log.SetError(err).Info(ctx, "can't save summary to redis")
		return entity.ResultSummary{}, err
	}

	return entity.ResultSummary{
		Result: listSummaryStock,
	}, nil
}

func updateSummary(recordedSummary *entity.OHLCSummary, newSummary entity.OHLCSummary) {

	if recordedSummary.PreviousPrice == 0 {
		recordedSummary.PreviousPrice = newSummary.PreviousPrice
	}

	if recordedSummary.LowestPrice == 0 {
		recordedSummary.LowestPrice = newSummary.LowestPrice
	}

	if newSummary.HighestPrice > recordedSummary.HighestPrice {
		recordedSummary.HighestPrice = newSummary.HighestPrice
	}

	if newSummary.LowestPrice > 0 && (newSummary.LowestPrice < recordedSummary.LowestPrice) {
		recordedSummary.LowestPrice = newSummary.LowestPrice
	}

	if newSummary.ClosePrice != 0 {
		recordedSummary.ClosePrice = newSummary.ClosePrice
	}

	recordedSummary.Value += newSummary.Value
	recordedSummary.Volume += newSummary.Volume
	recordedSummary.TotalTrans += newSummary.TotalTrans

	if recordedSummary.Value > 0 {
		recordedSummary.Average = float64(recordedSummary.Value) / float64(recordedSummary.Volume)
		recordedSummary.AverageRound = math.Round(recordedSummary.Average)
	}
}

func generateListSummaryStock(transactions []entity.OHLCTransaction) entity.ResultSummary {
	listSummaryStock := make(map[string]entity.OHLCSummary)
	listStock := make(map[string]int)

	for _, val := range transactions {
		if _, ok := listStock[val.StockCode]; !ok {
			listStock[val.StockCode] = 0
			listSummaryStock[val.StockCode] = entity.OHLCSummary{}
		}
		listStock[val.StockCode] += 1
	}

	for _, val := range transactions {
		var summary entity.OHLCSummary
		summary = listSummaryStock[val.StockCode]
		if val.Quantity == 0 {
			summary.PreviousPrice = val.Price
		}

		if summary.HighestPrice == 0 {
			summary.HighestPrice = val.Price
		}

		if summary.LowestPrice == 0 {
			summary.LowestPrice = val.Price
		}

		if valType := strings.ToLower(strings.TrimSpace(val.Type)); valType == "e" || valType == "p" {
			calculateAccountableType(&summary, val)
		}

		if summary.TotalTrans >= listStock[val.Type] {
			if summary.Volume > 0 {
				summary.Average = float64(summary.Value) / float64(summary.Volume)
			}
		}

		listSummaryStock[val.StockCode] = summary
	}

	return entity.ResultSummary{
		Result: listSummaryStock,
	}
}

func calculateAccountableType(summary *entity.OHLCSummary, transaction entity.OHLCTransaction) {
	switch {
	case summary.TotalTrans < 1:
		summary.HighestPrice = transaction.Price
		summary.LowestPrice = transaction.Price
	case transaction.Price > summary.HighestPrice:
		summary.HighestPrice = transaction.Price
	case transaction.Price > 0 && (transaction.Price < summary.LowestPrice):
		summary.LowestPrice = transaction.Price
	}

	//if {
	//	summary.HighestPrice = transaction.Price
	//	summary.LowestPrice = transaction.Price
	//}

	//if transaction.Price > summary.HighestPrice {
	//	summary.HighestPrice = transaction.Price
	//}
	//if transaction.Price > 0 && (transaction.Price < summary.LowestPrice) {
	//	summary.LowestPrice = transaction.Price
	//}

	summary.ClosePrice = transaction.Price
	summary.Volume += int64(transaction.Quantity)
	summary.Value += int64(transaction.Quantity * transaction.Price)
	summary.TotalTrans++
}
