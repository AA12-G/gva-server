package service

import (
	"context"
	"gva/internal/domain/entity"

	"gorm.io/gorm"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

// GetRolePermissions 获取角色的权限列表
func (s *RoleService) GetRolePermissions(ctx context.Context, roleID uint) ([]entity.Permission, error) {
	var role entity.Role
	err := s.db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

// AssignPermissions 为角色分配权限
func (s *RoleService) AssignPermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 先清除原有权限
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
			return err
		}

		// 分配新权限
		if len(permissionIDs) > 0 {
			var rolePermissions []map[string]interface{}
			for _, pid := range permissionIDs {
				rolePermissions = append(rolePermissions, map[string]interface{}{
					"role_id":       roleID,
					"permission_id": pid,
				})
			}
			return tx.Table("role_permissions").Create(rolePermissions).Error
		}
		return nil
	})
}

// TODO: 添加角色相关的业务逻辑
