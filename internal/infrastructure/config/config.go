package config

import (
	"gva/internal/infrastructure/redis"
	"os"

	"gva/internal/pkg/config"

	"gopkg.in/yaml.v3"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `yaml:"secret"`
	Expire int    `yaml:"expire"` // token过期时间（小时）
}

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	MySQL  config.MySQLConfig `yaml:"mysql"`
	Redis  redis.RedisConfig  `yaml:"redis"`
	JWT    JWTConfig          `yaml:"jwt"`
	Export struct {
		Dir string `yaml:"dir"`
	} `yaml:"export"`
}

func Load() (*Config, error) {
	var config Config
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
