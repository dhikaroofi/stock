package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

const OHLCSummaryKey = "OHLC#%s#%s"

type Task interface {
	SaveOHLCSummary(ctx context.Context, key string, summary entity.OHLCSummary) error
	SaveOHLCSummaryBatch(ctx context.Context, date string, summary map[string]entity.OHLCSummary) error
	GetOHLCSummary(ctx context.Context, stockCode, date string) (entity.OHLCSummary, error)
	GenerateKeyOHLCSummary(stockCode, date string) string
}

type redisCache struct {
	client *redis.Client
	ttl    time.Duration
	mutex  sync.Mutex
}

func New(client *redis.Client, ttl time.Duration) Task {
	return &redisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *redisCache) SaveOHLCSummary(ctx context.Context, key string, summary entity.OHLCSummary) error {

	payload, err := json.Marshal(summary)
	if err != nil {
		return err
	}

	pipe := c.client.TxPipeline()
	if exec := pipe.LPush(ctx, key, payload); exec.Err() != nil {
		pipe.Discard()
		return exec.Err()
	}

	if _, err := pipe.Exec(ctx); err != nil {
		pipe.Discard()
		return err
	}

	return nil
}

func (c *redisCache) SaveOHLCSummaryBatch(ctx context.Context, date string, summary map[string]entity.OHLCSummary) error {
	_, err := c.client.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for stockCode, val := range summary {
			key := c.GenerateKeyOHLCSummary(stockCode, date)

			payload := make(map[string]interface{})
			payload["previous_price"] = val.PreviousPrice
			payload["highest_price"] = val.HighestPrice
			payload["lowest_price"] = val.LowestPrice
			payload["close_price"] = val.ClosePrice
			payload["volume"] = val.Volume
			payload["value"] = val.Value
			payload["average"] = val.Average
			payload["average_round"] = val.AverageRound
			payload["total_trans"] = val.TotalTrans

			err := c.client.HSet(ctx, key, payload).Err()
			if err != nil {
				return err
			}

		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *redisCache) GetOHLCSummary(ctx context.Context, stockCode, date string) (entity.OHLCSummary, error) {
	exec, err := c.client.HGetAll(ctx, c.GenerateKeyOHLCSummary(stockCode, date)).Result()
	if err != nil {
		return entity.OHLCSummary{}, err
	}

	summary := convertSummaryToStruct(exec)

	return summary, nil
}

func (c *redisCache) GenerateKeyOHLCSummary(stockCode, date string) string {
	return fmt.Sprintf(OHLCSummaryKey, date, stockCode)
}

func convertSummaryToStruct(payload map[string]string) entity.OHLCSummary {
	var summary entity.OHLCSummary
	summary.Date = payload["date"]
	summary.StockCode = payload["stock_code"]
	summary.PreviousPrice, _ = strconv.Atoi(payload["previous_price"])
	summary.HighestPrice, _ = strconv.Atoi(payload["highest_price"])
	summary.LowestPrice, _ = strconv.Atoi(payload["lowest_price"])
	summary.ClosePrice, _ = strconv.Atoi(payload["close_price"])
	summary.Volume, _ = strconv.ParseInt(payload["volume"], 10, 64)
	summary.Value, _ = strconv.ParseInt(payload["value"], 10, 64)
	summary.Average, _ = strconv.ParseFloat(payload["average"], 64)
	summary.AverageRound, _ = strconv.ParseFloat(payload["average_round"], 64)
	summary.TotalTrans, _ = strconv.Atoi(payload["total_trans"])
	return summary
}
