package router

import (
	"gva/internal/domain/service"
	"gva/internal/infrastructure/repository"
	"gva/internal/interfaces/handler"
	"gva/internal/interfaces/middleware"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, rdb *redis.Client, userService *service.UserService) *gin.Engine {
	r := gin.Default()

	// 添加路由日志
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

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

	// 初始化处理器
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
		userBase := authorized.Group("/user")
		{
			userBase.GET("/info", userHandler.GetUserInfo)
			userBase.GET("/role", userHandler.GetUserRole)
			userBase.PUT("/profile", userHandler.UpdateProfile)
			userBase.POST("/reset-password", userHandler.ResetPassword)
			userBase.POST("/avatar", userHandler.UploadAvatar)
		}

		// 用户管理相关（需要用户管理权限）
		userManage := authorized.Group("/users")
		userManage.Use(middleware.CheckPermission("system:user"))
		{
			userManage.GET("", userHandler.ListUsers)
			userManage.GET("/:id/profile", userHandler.GetUserProfile)  // 获取指定用户信息
			userManage.PUT("/:id/profile", userHandler.UpdateUser)      // 更新指定用户信息
			userManage.PUT("/:id/status", userHandler.UpdateUserStatus) // 修改指定用户状态
			userManage.DELETE("/:id", userHandler.DeleteUser)           // 删除指定用户
			userManage.GET("/export", userHandler.ExportUsers)          // 导出用户数据
			userManage.POST("/import", userHandler.ImportUsers)         // 导入用户数据
			userManage.PUT("/:id/restore", userHandler.RestoreUser)     // 恢复已删除的用户
			userManage.POST("", userHandler.CreateUser)                 // 创建用户
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
		roleManage := authorized.Group("/roles")
		roleManage.Use(middleware.CheckPermission("system:role"))
		{
			roleManage.GET("", roleHandler.GetRoleList)                        // 获取角色列表
			roleManage.GET("/:id", roleHandler.GetRole)                        // 获取单个角色
			roleManage.POST("", roleHandler.CreateRole)                        // 创建角色
			roleManage.PUT("/:id", roleHandler.UpdateRole)                     // 更新角色
			roleManage.DELETE("/:id", roleHandler.DeleteRole)                  // 删除角色
			roleManage.GET("/:id/permissions", roleHandler.GetPermissions)     // 获取角色权限
			roleManage.POST("/:id/permissions", roleHandler.AssignPermissions) // 分配权限
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
