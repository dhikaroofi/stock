package command

import (
	"context"
	"encoding/json"
	"os"

	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/internal/core"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/logger"
)

// Task is an interface contain task for command
type Task interface {
	Run()
}

type command struct {
	config *config.Entity
	core   *core.Container
}

// New is constructor
func New(config *config.Entity, core *core.Container) Task {
	return &command{
		config: config,
		core:   core,
	}
}

// Run is for running challenge part 2
func (c command) Run() {
	var (
		newRecords  []entity.NewRecords
		indexMember []entity.IndexMember
		log         = logger.Call()
		ctx         = context.Background()
	)

	ctx = context.WithValue(ctx, logger.LogKeyType(logger.LogKey), log)

	if err := jsonParser(c.config.ChallengePart2.IndexMember, &indexMember); err != nil {
		log.SetError(err).TDR(ctx, "index-member.json")
		return
	}

	if err := jsonParser(c.config.ChallengePart2.NewRecords, &newRecords); err != nil {
		log.SetError(err).TDR(ctx, "new-record.json")
		return
	}

	if _, err := c.core.UseCase.ChallengePart2.OHLCUpdater(ctx, entity.OHLCIndexMember{
		Stock:       c.config.ChallengePart2.Stock,
		NewRecords:  newRecords,
		IndexMember: indexMember,
	}); err != nil {
		log.SetError(err).TDR(ctx, "error")
		return
	}

	log.TDR(ctx, "success")
}

func jsonParser(path string, payload interface{}) error {
	//nolint: gosec // this variable is get from config so it can be controlled
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&payload)
	if err != nil {
		return err
	}
	return nil
}
