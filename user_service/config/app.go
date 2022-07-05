package config

import (
	"github.com/pavel/user_service/pkg/logger"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	Server server
	logger *logger.Logger
}

var (
	instance *Config
	once     sync.Once
)

func InitConfig(logger *logger.Logger) (error, *Config) {
	var err error
	once.Do(func() {
		logger.Info("Init config")
		instance = &Config{logger: logger}
		err = instance.unmarshal()
	})
	if err != nil {
		logger.Fatal(err.Error())
		return err, nil
	}

	return nil, instance
}

func (cfg *Config) unmarshal() error {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		cfg.logger.Fatal(err.Error())
		return err
	}

	err = viper.Unmarshal(&cfg.Server)
	if err != nil {
		cfg.logger.Fatal(err.Error())
		return err
	}

	return nil
}
