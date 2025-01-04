package config

import (
	"gva/internal/infrastructure/redis"
	"os"

	"gva/internal/pkg/config"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	MySQL config.MySQLConfig `yaml:"mysql"`
	Redis redis.RedisConfig  `yaml:"redis"`
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
