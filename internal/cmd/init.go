package cmd

import (
	"fmt"
	"github.com/dhikaroofi/stock.git/internal/adapters/driving/datastreamer"
	"time"

	"github.com/dhikaroofi/stock.git/internal/adapters/driven/cache"
	"github.com/dhikaroofi/stock.git/internal/adapters/driving/grpc"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/internal/core"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	redis2 "github.com/dhikaroofi/stock.git/pkg/redis"
)

func Init(conf *config.Entity, existSignalch chan bool) {
	redisClient := redis2.NewRedis(redis2.Config{
		Host:     conf.Redis.Host,
		Password: conf.Redis.Password,
		DBIndex:  conf.Redis.Database,
	})

	cacheAdapter := cache.New(redisClient, time.Duration(conf.Redis.TTL)*time.Second)
	coreContainer := core.New(conf, &core.DrivenAdapter{
		Cache: cacheAdapter,
	})

	grpcServer := grpc.New(conf.GrpcPort, coreContainer)
	grpcServer.ListenAndServe(existSignalch)

	streamer := datastreamer.New(conf.DataStreamer.Path, coreContainer)
	streamer.ListenAndServe(existSignalch)

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
