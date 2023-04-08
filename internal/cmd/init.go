package cmd

import (
	data_streamer "github.com/dhikaroofi/stock.git/internal/adapters/driving/data-streamer"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/internal/core"
	"github.com/dhikaroofi/stock.git/pkg/logger"
)

func Init(conf *config.Entity, existSignalch chan bool) {

	var (
	//	readDB, readSqlDB   = gorm.Init(&conf.Database.Reader)
	//	writeDB, writeSqlDB = gorm.Init(&conf.Database.Writer)
	//	tdrLog              = zapLogger.InitZap(zapLogger.TDRlog, 1)
	//	sysLog              = zapLogger.InitZap(zapLogger.SYSlog, 1)
	//	firebase            = firebase2.Init(conf.Notification.ConfigPath)
	)
	//
	//repository := gorm_repo.NewGormRepo(readDB, writeDB)
	//notification := notification2.NewFirebaseNotification(firebase)
	//logs := logger.NewZapLogger(tdrLog, sysLog)
	//fileManager := storage.NewLocalFileManagement(conf.Storage.StoragePath,
	//	conf.Storage.Url, conf.Storage.MaxSize)

	coreContainer := core.New(conf, &core.DrivenAdapter{
		Cache: nil,
	})

	streamer := data_streamer.New(conf.DataStreamer.Path, coreContainer)
	streamer.ListenAndServe(existSignalch)
	//serverExitSignal := http.RunHttpServer(conf.App.Host, conf.App.Port, cont)

	go func() {
		<-existSignalch
		logger.SysInfo("disconnecting all dependent service")

		//if err := readSqlDB.Close(); err != nil {
		//	log.Error(err)
		//}
		//if err := writeSqlDB.Close(); err != nil {
		//	log.Error(err)
		//}
		logger.SysInfo("all the dependent services are disconnected")
		existSignalch <- true

	}()

}
