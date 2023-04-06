package redis

import (
	"github.com/redis/go-redis/v9"
)

func NewRedis(host, password string, dbIndex int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password, // no password set
		DB:       dbIndex,
	})
}
