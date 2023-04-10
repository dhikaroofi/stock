package config

import (
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"github.com/spf13/viper"
	"log"
)

const configFileType = "yaml"

type Entity struct {
	Rules struct {
		Stock []string
	}
	DataStreamer struct {
		Path string
	}
	Redis struct {
		Host     string
		Password string
		Database int
		TTL      int
	}
}

func LoadConfigFile(configName string) *Entity {
	viperConfig := viper.New()
	viperConfig.SetConfigName(configName)
	viperConfig.SetConfigType(configFileType)
	viperConfig.AddConfigPath(".")

	if err := viperConfig.ReadInConfig(); err != nil {
		log.Fatalf("failed reading config, %v", err)
	}

	config := new(Entity)
	if err := viperConfig.Unmarshal(&config); err != nil {
		log.Fatalf("failed parsing config, %v", err)
	}

	logger.SysInfo("config loaded")

	return config
}
