package core

import (
	cahce "github.com/dhikaroofi/stock.git/internal/adapters/driven/cache"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/internal/core/use-case/ohlc"
)

// Container is used for contain several use cases
type Container struct {
	UseCase usecase
}

type usecase struct {
	OHLC ohlc.Task
}

// DrivenAdapter is used for contain and parsing driven adapter from initiate process to container
type DrivenAdapter struct {
	Cache cahce.Task
}

// New this function is used for initiate the core
func New(conf *config.Entity, adapter *DrivenAdapter) *Container {
	var cont Container

	cont.UseCase.OHLC = ohlc.New(adapter.Cache)
	return &cont
}
