package router

import (
	"gva/internal/domain/service"
	"gva/internal/infrastructure/cache"
	"gva/internal/infrastructure/repository"
	"gva/internal/interfaces/handler"
	"gva/internal/interfaces/middleware"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 允许前端域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 配置静态文件服务
	r.Static("/uploads", "./uploads")

	// 初始化服务
	userRepo := repository.NewUserRepository(db)
	userCache := cache.NewUserCache(rdb)
	userService := service.NewUserService(userRepo, db, userCache)
	userHandler := handler.NewUserHandler(userService)

	logRepo := repository.NewOperationLogRepository(db)
	logService := service.NewOperationLogService(logRepo)
	logHandler := handler.NewOperationLogHandler(logService)

	// 添加操作日志中间件
	r.Use(middleware.OperationLog(logService))

	// 公开路由
	public := r.Group("/api/v1")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	authorized := r.Group("/api/v1")
	authorized.Use(middleware.JWTAuth())
	{
		// 用户相关
		authorized.PUT("/user/profile", userHandler.UpdateProfile)
		authorized.POST("/user/reset-password", userHandler.ResetPassword)
		authorized.POST("/user/avatar", userHandler.UploadAvatar)
		authorized.GET("/user/info", userHandler.GetUserInfo)

		// 用户管理
		authorized.GET("/users", userHandler.ListUsers)
		authorized.PUT("/users/:id/status", userHandler.UpdateUserStatus)
		authorized.DELETE("/users/:id", userHandler.DeleteUser)
		authorized.GET("/users/export", userHandler.ExportUsers)
		authorized.POST("/users/import", userHandler.ImportUsers)

		// 操作日志
		authorized.GET("/logs", logHandler.ListLogs)
	}

	return r
}
