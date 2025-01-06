package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"gva/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 打印完整的请求信息
		fmt.Printf("\n=== 请求信息 ===\n")
		fmt.Printf("Method: %s\n", c.Request.Method)
		fmt.Printf("Path: %s\n", c.Request.URL.Path)
		fmt.Printf("Headers: %+v\n", c.Request.Header)

		auth := c.GetHeader("Authorization")
		fmt.Printf("Authorization: %s\n", auth)

		// 如果是 OPTIONS 请求，直接放行
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			fmt.Printf("Token解析错误: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证信息"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
