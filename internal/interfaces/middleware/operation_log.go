package middleware

import (
	"bytes"
	"gva/internal/domain/entity"
	"gva/internal/domain/service"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// OperationLog 操作日志中间件
func OperationLog(logService *service.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 获取请求信息
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装 ResponseWriter 以获取响应内容
		w := &responseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		// 处理请求
		c.Next()

		// 获取用户ID
		userID, _ := c.Get("userID")

		// 创建日志记录
		log := &entity.OperationLog{
			UserID:    userID.(uint),
			IP:        c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Status:    c.Writer.Status(),
			Latency:   time.Since(startTime).Milliseconds(),
			UserAgent: c.Request.UserAgent(),
			Request:   string(requestBody),
			Response:  w.body.String(),
		}

		// 异步保存日志
		go logService.Create(c.Request.Context(), log)
	}
}
