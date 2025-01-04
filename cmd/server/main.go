package main

import (
	"log"

	"gva/internal/infrastructure/config"
	"gva/internal/infrastructure/database"
	"gva/internal/infrastructure/redis"
	"gva/internal/interfaces/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewMySQLConnection(cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 执行数据库迁移
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}
	defer sqlDB.Close()

	// 初始化Redis连接
	rdb, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	// 设置gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化路由
	r := router.InitRouter(db, rdb)

	// 启动服务器
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
