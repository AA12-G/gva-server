package handler

import (
	"fmt"
	"gva/internal/domain/entity"
	"gva/internal/domain/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService *service.PermissionService
}

func NewPermissionHandler(permissionService *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

// List 获取权限列表
func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.permissionService.GetAllPermissions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取权限列表失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": permissions,
	})
}

// Create 创建权限
func (h *PermissionHandler) Create(c *gin.Context) {
	var permission entity.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: 实现创建权限的逻辑
	c.JSON(http.StatusOK, gin.H{"message": "创建成功"})
}

// Update 更新权限
func (h *PermissionHandler) Update(c *gin.Context) {
	var permission entity.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: 实现更新权限的逻辑
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除权限
func (h *PermissionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	// 将字符串ID转换为uint
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的权限ID"})
		return
	}

	// TODO: 调用service层删除权限
	_ = uint(idUint) // 暂时使用空白标识符避免未使用错误

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
