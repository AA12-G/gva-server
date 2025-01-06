package service

import (
	"context"

	"gorm.io/gorm"

	"gva/internal/domain/entity"
)

type PermissionService struct {
	db *gorm.DB
}

func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// 获取用户权限
func (s *PermissionService) GetUserPermissions(ctx context.Context, userID uint) ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := s.db.Model(&entity.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Joins("JOIN users ON users.role_id = roles.id").
		Where("users.id = ? AND permissions.status = 1", userID).
		Find(&permissions).Error
	return permissions, err
}

// 检查用户是否有特定权限
func (s *PermissionService) HasPermission(ctx context.Context, userID uint, permissionCode string) bool {
	var count int64
	s.db.Model(&entity.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Joins("JOIN users ON users.role_id = roles.id").
		Where("users.id = ? AND permissions.code = ? AND permissions.status = 1", userID, permissionCode).
		Count(&count)
	return count > 0
}
