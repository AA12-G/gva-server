package middleware

import (
	"log"
	"net/http"

	"gva/internal/domain/service"

	"github.com/gin-gonic/gin"
)

func CheckPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("检查权限: %s, 路径: %s", permission, c.Request.URL.Path)

		// 从上下文获取用户ID
		userIDInterface, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}

		// 类型断言，确保安全转换
		userID, ok := userIDInterface.(uint)
		if !ok {
			log.Printf("用户ID类型错误: %T", userIDInterface)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
			c.Abort()
			return
		}

		permissionService, exists := c.Get("permissionService")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限服务未初始化"})
			c.Abort()
			return
		}

		hasPermission := permissionService.(*service.PermissionService).HasPermission(
			c.Request.Context(),
			userID,
			permission,
		)

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}
