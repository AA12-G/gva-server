package config

import (
	"gva/internal/pkg/config"

	"github.com/spf13/viper"
)

type Config struct {
	Server config.ServerConfig `mapstructure:"server"`
	MySQL  config.MySQLConfig  `mapstructure:"mysql"`
	Redis  config.RedisConfig  `mapstructure:"redis"`
	JWT    config.JWTConfig    `mapstructure:"jwt"`
	Export config.ExportConfig `mapstructure:"export"`
}

func LoadConfig(file string) (*Config, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
