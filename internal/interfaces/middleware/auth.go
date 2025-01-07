package middleware

import (
	"github.com/gin-gonic/gin"
)

// 重命名为 AuthMiddleware 或者完全删除这个文件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ... 原来的代码 ...
	}
}
