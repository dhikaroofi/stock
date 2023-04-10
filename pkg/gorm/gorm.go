package gorm

import (
	"database/sql"
	"fmt"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(conf ConfigDB) (*gorm.DB, *sql.DB) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Pass,
		conf.Host, conf.Port, conf.DatabaseName)

	driverConfig := mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	})

	gormConfig := &gorm.Config{
		PrepareStmt: true,
	}

	provider, err := gorm.Open(driverConfig, gormConfig)
	if err != nil {
		logger.Fatal(fmt.Sprintf("model not connected: %s", err.Error()))
	}

	sqlDB, err := provider.DB()
	if err != nil {
		logger.Fatal(fmt.Sprintf("model not connected: %s", err.Error()))
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Fatal(fmt.Sprintf("failed checking connection to database, err: %s", err.Error()))
	}

	sqlDB.SetMaxIdleConns(conf.PoolMaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.PoolMaxOpenConn)
	sqlDB.SetConnMaxLifetime(conf.PoolMaxConnLifetime)

	logger.SysInfo("gormDB : successfully connected to database")

	if conf.DebugLog {
		provider = provider.Debug()
	}

	return provider, sqlDB
}
