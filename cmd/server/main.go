package main

import (
	"log"

	"gva/internal/domain/cache"
	"gva/internal/domain/service"
	infraCache "gva/internal/infrastructure/cache"
	"gva/internal/infrastructure/config"
	"gva/internal/infrastructure/database"
	"gva/internal/infrastructure/redis"
	"gva/internal/infrastructure/repository"
	"gva/internal/interfaces/router"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库连接
	db, err := database.NewMySQLDB(&cfg.MySQL)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 自动迁移数据库
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化 Redis 客户端
	redisConfig := &redis.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}

	rdb, err := redis.NewRedisClient(redisConfig)
	if err != nil {
		log.Printf("Redis连接失败，将在无缓存模式下运行: %v", err)
	}

	// 初始化服务
	userRepo := repository.NewUserRepository(db)
	var userCache cache.UserCache
	if rdb != nil {
		userCache = infraCache.NewRedisUserCache(rdb)
	}
	userService := service.NewUserService(userRepo, db, userCache)

	// 初始化路由
	r := router.InitRouter(db, rdb, userService)

	// 启动服务器
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
