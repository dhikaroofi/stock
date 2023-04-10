package config

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dhikaroofi/stock.git/pkg/logger"
	"github.com/spf13/viper"
)

const configFileType = "yaml"

var (
	RootDir = ""
)

// Entity is for config value
type Entity struct {
	RootDir        string
	ChallengePart2 challengePart2
	DataStreamer   dataStreamer
	Redis          redis
}

type challengePart2 struct {
	Stock       []string
	IndexMember string
	NewRecords  string
}

type dataStreamer struct {
	Path string
}

type redis struct {
	Host     string
	Password string
	Database int
	TTL      int
}

// LoadConfigFile is function to load config from config.yaml
func LoadConfigFile(configName string) *Entity {
	viperConfig := viper.New()
	viperConfig.SetConfigName(configName)
	viperConfig.SetConfigType(configFileType)
	viperConfig.AddConfigPath(".")

	if err := viperConfig.ReadInConfig(); err != nil {
		logger.Fatal(fmt.Sprintf("failed reading config, %s", err.Error()))
	}

	config := new(Entity)
	if err := viperConfig.Unmarshal(&config); err != nil {
		logger.Fatal(fmt.Sprintf("failed parsing config, %s", err.Error()))
	}

	getRootDir()

	logger.SysInfo("config loaded")

	return config
}

func getRootDir() {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	out, err := cmd.Output()
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed get root dir, %s", err.Error()))
	}
	// Trim the newline character from the output
	RootDir = strings.TrimSpace(string(out))
}
