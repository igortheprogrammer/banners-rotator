package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Logger  LoggerConf  `yaml:"logger"`
	Api     ApiConf     `yaml:"api"`
	Storage StorageConf `yaml:"storage"`
	Rmq     RmqConf     `yaml:"rmq"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type ApiConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type StorageConf struct {
	Type             string `yaml:"type"`
	ConnectionString string `yaml:"connectionString"`
}

type RmqConf struct {
	Uri  string `yaml:"uri"`
	Name string `yaml:"name"`
}

var ErrUnreadableConfig = errors.New("unreadable config")

func init() {
	viper.SetDefault("logger.level", "info")
}

func NewAppConfig(path string) (AppConfig, error) {
	dir, file := filepath.Split(path)
	name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

	viper.SetConfigName(name)
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		return AppConfig{}, fmt.Errorf("%w: %v", ErrUnreadableConfig, err)
	}

	cfg := AppConfig{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return AppConfig{}, fmt.Errorf("config unmarshalling: %w", err)
	}

	return cfg, nil
}
