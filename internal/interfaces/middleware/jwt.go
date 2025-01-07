package middleware

import (
	"fmt"
	"gva/internal/pkg/jwt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("收到请求: %s %s", c.Request.Method, c.Request.URL.Path)

		// 获取 Authorization 头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		// 验证 token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("无效的token: %v", err)})
			c.Abort()
			return
		}

		// 将用户ID存入上下文（确保是 uint 类型）
		c.Set("userID", uint(claims.UserID))
		c.Next()
	}
}
