package router

import (
	"gva/internal/domain/service"
	"gva/internal/infrastructure/repository"
	"gva/internal/interfaces/handler"
	"gva/internal/interfaces/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	// 初始化处理器
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)
	userHandler := handler.NewUserHandler(userService)

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
		authorized.PUT("/user/profile", userHandler.UpdateProfile)
		authorized.POST("/user/reset-password", userHandler.ResetPassword)
	}

	return r
}
