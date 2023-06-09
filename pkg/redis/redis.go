package redis

import (
	"context"
	"fmt"

	"github.com/dhikaroofi/stock.git/pkg/logger"
	"github.com/redis/go-redis/v9"
)

// Config this struct is used for parsing config from outside
type Config struct {
	Host     string
	Password string
	DBIndex  int
}

// NewRedis is function for called redis client
func NewRedis(config Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Password, // no password set
		DB:       config.DBIndex,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.Fatal(fmt.Sprintf("redis is unreachable: %s", err.Error()))
	}

	return client
}
