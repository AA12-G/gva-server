package config

import (
	"fmt"
	"gva/internal/infrastructure/database"
	"gva/internal/infrastructure/redis"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
	MySQL database.MySQLConfig `yaml:"mysql"`
	Redis redis.RedisConfig    `yaml:"redis"`
}

func Load() (*Config, error) {
	var config Config

	// 打印当前工作目录，帮助调试
	pwd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", pwd)

	// 读取配置文件 - 使用相对于项目根目录的路径
	data, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("read config file error: %v, pwd: %s", err, pwd)
	}

	// 解析YAML并打印内容，帮助调试
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	fmt.Printf("Loaded config: %+v\n", config)

	return &config, nil
}
