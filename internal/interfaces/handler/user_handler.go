package handler

import (
	"fmt"
	"gva/internal/domain/service"
	"gva/internal/pkg/jwt"
	"gva/internal/pkg/upload"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	err := h.userService.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	user, err := h.userService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user":    user,
		"token":   token,
	})
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email" binding:"omitempty,email"`
		Phone    string `json:"phone" binding:"omitempty,numeric,len=11"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		var errMsg string
		switch {
		case err.Error() == "Key: 'Phone' Error:Field validation for 'Phone' failed on the 'len' tag":
			errMsg = "手机号必须是11位数字"
		case err.Error() == "Key: 'Phone' Error:Field validation for 'Phone' failed on the 'numeric' tag":
			errMsg = "手机号只能包含数字"
		case err.Error() == "Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag":
			errMsg = "邮箱格式不正确"
		default:
			errMsg = "参数错误"
		}
		fmt.Printf("验证错误: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	err := h.userService.UpdateProfile(c.Request.Context(), userID.(uint), req.Nickname, req.Email, req.Phone, req.Avatar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// ResetPassword 重置密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, _ := c.Get("userID")
	err := h.userService.ResetPassword(c.Request.Context(), userID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码修改成功"})
}

// UploadAvatar 上传头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	// 保存文件
	filePath, err := upload.SaveUploadedFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新用户头像
	err = h.userService.UpdateAvatar(c.Request.Context(), userID.(uint), filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"avatar":  filePath,
	})
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// ListUsers 获取用户列表
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req struct {
		Page     int    `form:"page" binding:"omitempty,min=1"`
		PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
		Keyword  string `form:"keyword"`
		Status   *int   `form:"status"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("验证错误: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10 // 默认每页10条
	}

	// 添加请求参数日志
	fmt.Printf("查询参数: page=%d, pageSize=%d, keyword=%s, status=%v\n",
		req.Page, req.PageSize, req.Keyword, req.Status)

	users, total, err := h.userService.ListUsers(c.Request.Context(), req.Page, req.PageSize, req.Keyword, req.Status)
	if err != nil {
		fmt.Printf("查询错误: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 添加结果日志
	fmt.Printf("查询结果: 总数=%d, 返回记录数=%d\n", total, len(users))

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}

// UpdateUserStatus 更新用户状态
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	var req struct {
		UserID uint `uri:"id" binding:"required"`
		Status int  `json:"status" binding:"required,oneof=0 1 2"` // 0:禁用 1:正常 2:待审核
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID无效"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态参数无效"})
		return
	}

	if err := h.userService.UpdateUserStatus(c.Request.Context(), req.UserID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}
