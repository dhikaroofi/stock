package cmd

import (
	"fmt"
	"time"

	"github.com/dhikaroofi/stock.git/internal/adapters/driven/cache"
	"github.com/dhikaroofi/stock.git/internal/adapters/driving/command"
	"github.com/dhikaroofi/stock.git/internal/adapters/driving/datastreamer"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/internal/core"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	redis2 "github.com/dhikaroofi/stock.git/pkg/redis"
)

// InitChallengePart1 comment part
func InitChallengePart1(conf *config.Entity, existSignalch chan bool) {
	redisClient := redis2.NewRedis(redis2.Config{
		Host:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DBIndex:  conf.Redis.Database,
	})

	cacheAdapter := cache.New(redisClient, time.Duration(conf.Redis.TTL)*time.Second)
	coreContainer := core.New(conf, &core.DrivenAdapter{
		Cache: cacheAdapter,
	})

	streamer := datastreamer.New(conf.DataStreamer.Path, coreContainer)
	streamer.ListenAndServe(existSignalch)

	logger.SysInfo("part 1 is on the way")
	go func() {
		<-existSignalch
		logger.SysInfo("disconnecting all dependent service")

		if err := redisClient.Close(); err != nil {
			logger.SysInfo(fmt.Sprintf("cannot close redis client: %s", err.Error()))
		}

		logger.SysInfo("all the dependent services are disconnected")
		existSignalch <- true
	}()
}

// InitChallengePart2 is used for initiate part 2 challenge
func InitChallengePart2(conf *config.Entity) {
	coreContainer := core.New(conf, &core.DrivenAdapter{
		Cache: nil,
	})

	logger.SysInfo("part 2 is on the way")
	cmd := command.New(conf, coreContainer)
	cmd.Run()
}
