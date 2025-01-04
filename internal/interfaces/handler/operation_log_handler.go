package handler

import (
	"gva/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OperationLogHandler struct {
	logService *service.OperationLogService
}

func NewOperationLogHandler(logService *service.OperationLogService) *OperationLogHandler {
	return &OperationLogHandler{logService: logService}
}

// ListLogs 获取操作日志列表
func (h *OperationLogHandler) ListLogs(c *gin.Context) {
	var req struct {
		Page     int `form:"page" binding:"omitempty,min=1"`
		PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	logs, total, err := h.logService.List(c.Request.Context(), req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}
