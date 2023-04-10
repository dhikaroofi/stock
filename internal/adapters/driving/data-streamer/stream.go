package data_streamer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhikaroofi/stock.git/internal/core"
	entity2 "github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"io"
	"os"
	"path/filepath"
)

const ext = ".ndjson"
const stop = false

type Task interface {
	ListenAndServe(exitSignal chan bool)
	Close()
}

type stream struct {
	path string
	core *core.Container
}

func New(path string, coreContainer *core.Container) Task {
	return &stream{
		path: path,
		core: coreContainer,
	}
}

func (s stream) ListenAndServe(exitSignal chan bool) {

	files, err := os.ReadDir(s.path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {

		// if in kafka this thing can be consume data on specific topic
		if filepath.Ext(file.Name()) == ext {
			var (
				logCaller = logger.Call()
				ctx       = context.Background()
			)

			ctx = context.WithValue(ctx, logger.LogKey, logCaller)

			filePath := filepath.Join(s.path, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				logCaller.SetError(err).TDR(ctx, fmt.Sprintf("error on read data"))
				continue
			}

			// parsing raw data into slice of transactions
			transactions, err := s.Parser(bytes.NewReader(content))
			if err != nil {
				logCaller.SetError(err).TDR(ctx, fmt.Sprintf("error on parser data"))
				continue
			}

			// calculating data and get summary
			if err := s.core.UseCase.OHLC.Calculate(ctx, file.Name(), transactions); err != nil {
				logCaller.SetError(err).TDR(ctx, "")
				continue
			}

			logCaller.TDR(ctx, "success")
		}

		if exitSignal != nil && len(exitSignal) > 0 {
			<-exitSignal
			return
		}

	}

	go func() {
		<-exitSignal
		s.Close()
		logger.SysInfo("data streamer is shutdown")
		exitSignal <- true
	}()

}

func (s stream) Parser(items io.Reader) ([]entity2.OHLCTransaction, error) {
	var transactions []entity2.OHLCTransaction

	scanner := bufio.NewScanner(items)

	for scanner.Scan() {

		var entity entity2.OHLCTransactionReq

		err := json.Unmarshal(scanner.Bytes(), &entity)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, entity.Translate())
	}

	if err := scanner.Err(); err != nil {
		return nil, scanner.Err()
	}

	return transactions, nil
}

func (s stream) Close() {}
