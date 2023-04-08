package ohlc

import (
	"context"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	dlog "log"
)

type Task interface {
	Calculate(ctx context.Context, transactions []entity.OHLCTransaction) error
	GetSummary(ctx context.Context, id string) error
}

type ohlcUseCase struct{}

func New() Task {
	return &ohlcUseCase{}
}

func (uc ohlcUseCase) Calculate(ctx context.Context, transactions []entity.OHLCTransaction) error {
	//var listStockTran = map[string]entity.OHLCSummary{}

	log := logger.ExtractLogFromContext(ctx)

	listStock := make(map[string]int)
	go func() {
		for _, val := range transactions {
			if _, ok := listStock[val.StockCode]; ok {
				listStock[val.StockCode] = 0
			}
			listStock[val.StockCode] = listStock[val.StockCode] + 1
		}
	}()

	dlog.Println(listStock)

	//listSummaryStock := make(map[string]entity.OHLCSummary, len(listStock))
	//for _, val := range transactions {
	//	var summary entity.OHLCSummary
	//	if _, ok := listSummaryStock[val.StockCode]; ok {
	//		summary = listSummaryStock[val.StockCode]
	//	}
	//
	//	if val.Quantity == 0 {
	//		summary.PreviousPrice = val.Price
	//	}
	//
	//	if valType := strings.ToLower(strings.TrimSpace(val.Type)); valType == "e" || valType == "p" {
	//		summary.Volume = summary.Volume + val.Quantity
	//		summary.Value = summary.Value + (val.Quantity * val.Price)
	//	}
	//
	//	listSummaryStock[val.StockCode] = summary
	//
	//}

	log.Info(ctx, "calculation is success")
	return nil
}

func (uc ohlcUseCase) GetSummary(ctx context.Context, id string) error {
	return nil
}
