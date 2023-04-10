package datastreamer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/dhikaroofi/stock.git/internal/core"
	entity2 "github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/logger"
)

const ext = ".ndjson"

// Task interface
type Task interface {
	ListenAndServe(exitSignal chan bool)
}

type stream struct {
	path string
	core *core.Container
}

// New is used for initiate data streamer from file or it can be kafka
func New(path string, coreContainer *core.Container) Task {
	return &stream{
		path: path,
		core: coreContainer,
	}
}

func (s stream) ListenAndServe(exitSignal chan bool) {
	files, err := os.ReadDir(s.path)
	if err != nil {
		logger.Fatal(fmt.Sprintf("invalid directory:%s", err.Error()))
	}

	mu := sync.Mutex{}
	for _, file := range files {
		// if in kafka this thing can be consume data on specific topic
		if filepath.Ext(file.Name()) == ext {
			mu.Lock()
			var (
				logCaller = logger.Call()
				ctx       = context.Background()
				fileName  = file.Name()
			)

			content, err := s.readJSONFile(fileName)
			if err != nil {
				logCaller.SetError(err).TDR(ctx, "error on read json file")
				continue
			}

			ctx = context.WithValue(ctx, logger.LogKeyType(logger.LogKey), logCaller)
			if err := s.handleData(ctx, fileName, content); err != nil {
				logCaller.SetError(err).TDR(ctx, "error on read data")
				continue
			}
			logCaller.TDR(ctx, fmt.Sprintf("success: %s", fileName))
			mu.Unlock()
		}

		//nolint: gosimple // this conditional is for handling gracefully shutdown to stop looping
		if exitSignal != nil && len(exitSignal) > 0 {
			exitSignal <- true
			return
		}
	}

	go func() {
		<-exitSignal
		logger.SysInfo("data streamer is shutdown")
		exitSignal <- true
	}()
}

func (s stream) readJSONFile(fileName string) ([]byte, error) {
	filePath := filepath.Join(s.path, fileName)

	//nolint: gosec // this variable is get from config so it can be controlled
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (s stream) handleData(ctx context.Context, fileName string, content []byte) error {
	// parsing raw data into slice of transactions
	transactions, err := s.Parser(bytes.NewReader(content))
	if err != nil {
		return err
	}

	// calculating data and get summary
	if _, err := s.core.UseCase.OHLC.Calculate(ctx, entity2.ListOHLCTransactions{
		Info: fileName,
		List: transactions,
	}); err != nil {
		return err
	}
	return nil
}

func (stream) Parser(items io.Reader) ([]entity2.OHLCTransaction, error) {
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
