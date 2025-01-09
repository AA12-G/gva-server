package service

import (
	"context"
	"errors"
	"fmt"
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

// GetAllRoles 获取所有角色列表
func (s *RoleService) GetAllRoles(ctx context.Context) ([]entity.Role, error) {
	var roles []entity.Role

	// 查询所有未删除的角色，并按sort和id排序
	err := s.db.Model(&entity.Role{}).
		Where("deleted_at IS NULL").
		Order("id ASC").
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("查询角色列表失败: %v", err)
	}

	return roles, nil
}

// GetRoleByID 根据ID获取角色信息
func (s *RoleService) GetRoleByID(ctx context.Context, id uint) (*entity.Role, error) {
	var role entity.Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("角色不存在")
		}
		return nil, fmt.Errorf("查询角色失败: %v", err)
	}
	return &role, nil
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(ctx context.Context, role *entity.Role) error {
	return s.db.Create(role).Error
}

// UpdateRole 更新角色信息
func (s *RoleService) UpdateRole(ctx context.Context, id uint, role *entity.Role) error {
	// 检查角色是否存在
	var existingRole entity.Role
	if err := s.db.First(&existingRole, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %v", err)
	}

	// 如果修改了角色编码，检查新编码是否已存在
	if role.Code != existingRole.Code {
		var count int64
		if err := s.db.Model(&entity.Role{}).Where("code = ? AND id != ?", role.Code, id).Count(&count).Error; err != nil {
			return fmt.Errorf("检查角色编码失败: %v", err)
		}
		if count > 0 {
			return fmt.Errorf("角色编码 %s 已存在", role.Code)
		}
	}

	// 如果未提供状态，默认设置为1
	if role.Status == 0 {
		role.Status = 1
	}

	// 更新角色信息
	updates := map[string]interface{}{
		"name":        role.Name,
		"code":        role.Code,
		"description": role.Description,
		"sort":        role.Sort,
		"status":      role.Status,
	}

	if err := s.db.Model(&existingRole).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新角色失败: %v", err)
	}

	return nil
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(ctx context.Context, id uint) error {
	return s.db.Delete(&entity.Role{}, id).Error
}

// TODO: 添加角色相关的业务逻辑
