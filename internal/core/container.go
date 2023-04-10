package core

import (
	cahce "github.com/dhikaroofi/stock.git/internal/adapters/driven/cache"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/internal/core/use-case/ohlc"
)

type Container struct {
	UseCase struct {
		OHLC ohlc.Task
	}
}

type DrivenAdapter struct {
	Cache cahce.Task
}

func New(conf *config.Entity, adapter *DrivenAdapter) *Container {
	var cont Container

	cont.UseCase.OHLC = ohlc.New(adapter.Cache)

	return &cont
}
