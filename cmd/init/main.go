package main

import (
	"gva/internal/infrastructure/config"
	"gva/internal/infrastructure/database"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := database.InitDB(&cfg.MySQL)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 强制执行初始化
	if err := database.InitData(db); err != nil {
		log.Fatalf("初始化数据失败: %v", err)
	}

	log.Println("数据初始化完成")
}
