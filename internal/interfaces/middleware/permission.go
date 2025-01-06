package middleware

import (
	"fmt"
	"net/http"

	"gva/internal/domain/service"

	"github.com/gin-gonic/gin"
)

func CheckPermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("\n=== 权限检查 ===\n")
		fmt.Printf("检查权限代码: %s\n", permissionCode)

		userID, exists := c.Get("userID")
		if !exists {
			fmt.Printf("未找到用户ID\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
			c.Abort()
			return
		}
		fmt.Printf("用户ID: %v\n", userID)

		permissionService, exists := c.Get("permissionService")
		if !exists {
			fmt.Printf("未找到权限服务\n")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限服务未初始化"})
			c.Abort()
			return
		}

		hasPermission := permissionService.(*service.PermissionService).HasPermission(
			c.Request.Context(),
			userID.(uint),
			permissionCode,
		)
		fmt.Printf("权限检查结果: %v\n", hasPermission)

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		c.Next()
	}
}
