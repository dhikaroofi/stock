package ohlc

import (
	"context"
	"testing"

	"github.com/dhikaroofi/stock.git/internal/adapters/driven/cache/cachefakes"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	transactions := []entity.OHLCTransaction{
		{StockCode: "BBCA", Type: "A", Price: 8000, Quantity: 0},
		{StockCode: "BBCA", Type: "P", Price: 8050, Quantity: 100},
		{StockCode: "BBCA", Type: "P", Price: 7950, Quantity: 500},
		{StockCode: "BBCA", Type: "A", Price: 8150, Quantity: 200},
		{StockCode: "BBCA", Type: "E", Price: 8100, Quantity: 300},
		{StockCode: "BBCA", Type: "A", Price: 8200, Quantity: 100},
	}

	expectedReulst := entity.OHLCSummary{
		Average:       8011.111111111111,
		AverageRound:  8011,
		Volume:        900,
		Value:         7210000,
		PreviousPrice: 8000,
		HighestPrice:  8100,
		LowestPrice:   7950,
		ClosePrice:    8100,
		TotalTrans:    3,
	}

	mockCache := cachefakes.FakeTask{}
	uc := New(&mockCache)

	ctx := context.Background()
	result, err := uc.Calculate(ctx, entity.ListOHLCTransactions{
		Info: "2021-01-01",
		List: transactions,
	})
	summary := result.Result

	assert.NoError(t, err)
	if val, ok := summary["BBCA"]; ok {
		assert.Equal(t, expectedReulst, val)
	} else {
		t.Errorf("expected map is filled with BBCA ")
	}
}
