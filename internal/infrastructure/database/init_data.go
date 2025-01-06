package database

import (
	"fmt"
	"gva/internal/domain/entity"
	"gva/internal/pkg/utils"
	"log"

	"gorm.io/gorm"
)

// InitData 初始化基础数据
func InitData(db *gorm.DB) error {
	log.Println("开始初始化基础数据...")

	// 1. 创建基础角色
	roles := []entity.Role{
		{
			Name:        "超级管理员",
			Code:        "super_admin",
			Status:      1,
			Sort:        1,
			Description: "系统超级管理员，拥有所有权限",
		},
		{
			Name:        "管理员",
			Code:        "admin",
			Status:      1,
			Sort:        2,
			Description: "系统管理员，拥有大部分权限",
		},
		{
			Name:        "普通用户",
			Code:        "user",
			Status:      1,
			Sort:        3,
			Description: "普通用户，只有基本权限",
		},
	}

	// 2. 创建基础权限
	permissions := []entity.Permission{
		// 用户管理权限
		{
			Name:        "用户管理",
			Code:        "system:user",
			Type:        "menu",
			Status:      1,
			Description: "用户管理菜单",
		},
		{
			Name:        "查看用户",
			Code:        "system:user:list",
			Type:        "button",
			Status:      1,
			Description: "查看用户列表",
		},
		{
			Name:        "创建用户",
			Code:        "system:user:create",
			Type:        "button",
			Status:      1,
			Description: "创建新用户",
		},
		{
			Name:        "更新用户",
			Code:        "system:user:update",
			Type:        "button",
			Status:      1,
			Description: "更新用户信息",
		},
		{
			Name:        "删除用户",
			Code:        "system:user:delete",
			Type:        "button",
			Status:      1,
			Description: "删除用户",
		},
		{
			Name:        "导出用户",
			Code:        "system:user:export",
			Type:        "button",
			Status:      1,
			Description: "导出用户数据",
		},
		{
			Name:        "导入用户",
			Code:        "system:user:import",
			Type:        "button",
			Status:      1,
			Description: "导入用户数据",
		},

		// 权限管理权限
		{
			Name:        "权限管理",
			Code:        "system:permission",
			Type:        "menu",
			Status:      1,
			Description: "权限管理菜单",
		},
		{
			Name:        "查看权限",
			Code:        "system:permission:list",
			Type:        "button",
			Status:      1,
			Description: "查看权限列表",
		},
		{
			Name:        "创建权限",
			Code:        "system:permission:create",
			Type:        "button",
			Status:      1,
			Description: "创建新权限",
		},
		{
			Name:        "更新权限",
			Code:        "system:permission:update",
			Type:        "button",
			Status:      1,
			Description: "更新权限信息",
		},
		{
			Name:        "删除权限",
			Code:        "system:permission:delete",
			Type:        "button",
			Status:      1,
			Description: "删除权限",
		},

		// 角色管理权限
		{
			Name:        "角色管理",
			Code:        "system:role",
			Type:        "menu",
			Status:      1,
			Description: "角色管理菜单",
		},
		{
			Name:        "查看角色",
			Code:        "system:role:list",
			Type:        "button",
			Status:      1,
			Description: "查看角色列表",
		},
		{
			Name:        "创建角色",
			Code:        "system:role:create",
			Type:        "button",
			Status:      1,
			Description: "创建新角色",
		},
		{
			Name:        "更新角色",
			Code:        "system:role:update",
			Type:        "button",
			Status:      1,
			Description: "更新角色信息",
		},
		{
			Name:        "删除角色",
			Code:        "system:role:delete",
			Type:        "button",
			Status:      1,
			Description: "删除角色",
		},
		{
			Name:        "分配权限",
			Code:        "system:role:assign",
			Type:        "button",
			Status:      1,
			Description: "为角色分配权限",
		},

		// 日志管理权限
		{
			Name:        "日志管理",
			Code:        "system:log",
			Type:        "menu",
			Status:      1,
			Description: "日志管理菜单",
		},
		{
			Name:        "查看日志",
			Code:        "system:log:list",
			Type:        "button",
			Status:      1,
			Description: "查看操作日志",
		},
	}

	// 3. 创建默认管理员用户
	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	adminUser := entity.User{
		Username: "admin",
		Password: hashedPassword, // 使用加密后的密码
		Nickname: "超级管理员",
		Email:    "admin@example.com",
		Status:   1,
		RoleID:   1, // 超级管理员角色
	}

	// 4. 执行数据初始化
	return db.Transaction(func(tx *gorm.DB) error {
		log.Println("创建角色...")
		if err := tx.Create(&roles).Error; err != nil {
			return fmt.Errorf("创建角色失败: %v", err)
		}

		log.Println("创建权限...")
		if err := tx.Create(&permissions).Error; err != nil {
			return fmt.Errorf("创建权限失败: %v", err)
		}

		log.Println("为超级管理员分配权限...")
		if err := tx.Model(&roles[0]).Association("Permissions").Replace(&permissions); err != nil {
			return fmt.Errorf("分配权限失败: %v", err)
		}

		log.Println("创建管理员用户...")
		if err := tx.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("创建管理员用户失败: %v", err)
		}

		log.Println("数据初始化完成")
		return nil
	})
}
