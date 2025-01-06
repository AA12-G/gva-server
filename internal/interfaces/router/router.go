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
		AllowOrigins:     []string{"*"}, // 允许所有来源，生产环境应该设置具体域名
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
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

	// 初始化权限相关服务和处理器
	permissionService := service.NewPermissionService(db)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	roleService := service.NewRoleService(db)
	roleHandler := handler.NewRoleHandler(roleService)

	// 将权限服务添加到全局上下文
	r.Use(func(c *gin.Context) {
		c.Set("permissionService", permissionService)
		c.Next()
	})

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
		// 基础功能（不需要额外权限）
		authorized.GET("/user/info", userHandler.GetUserInfo)
		authorized.PUT("/user/profile", userHandler.UpdateProfile)
		authorized.POST("/user/reset-password", userHandler.ResetPassword)
		authorized.POST("/user/avatar", userHandler.UploadAvatar)

		// 用户管理相关（需要用户管理权限）
		userManage := authorized.Group("")
		userManage.Use(middleware.CheckPermission("system:user"))
		{
			userManage.GET("/users", userHandler.ListUsers)
			userManage.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			userManage.DELETE("/users/:id", userHandler.DeleteUser)
			userManage.GET("/users/export", userHandler.ExportUsers)
			userManage.POST("/users/import", userHandler.ImportUsers)
		}

		// 权限管理相关（需要权限管理权限）
		permManage := authorized.Group("")
		permManage.Use(middleware.CheckPermission("system:permission"))
		{
			permManage.GET("/permissions", permissionHandler.List)
			permManage.POST("/permissions", permissionHandler.Create)
			permManage.PUT("/permissions/:id", permissionHandler.Update)
			permManage.DELETE("/permissions/:id", permissionHandler.Delete)
		}

		// 角色权限管理
		roleManage := authorized.Group("")
		roleManage.Use(middleware.CheckPermission("system:role"))
		{
			roleManage.GET("/roles/:id/permissions", roleHandler.GetPermissions)
			roleManage.POST("/roles/:id/permissions", roleHandler.AssignPermissions)
		}

		// 日志管理（需要日志查看权限）
		logManage := authorized.Group("")
		logManage.Use(middleware.CheckPermission("system:log"))
		{
			logManage.GET("/logs", logHandler.ListLogs)
		}
	}

	return r
}
