package database

import (
	"gva/internal/domain/entity"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	log.Println("开始迁移数据库...")
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
	)
	if err != nil {
		return err
	}

	// 创建默认角色
	var count int64
	db.Model(&entity.Role{}).Count(&count)
	if count == 0 {
		defaultRole := &entity.Role{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户角色",
			Status:      1,
		}
		if err := db.Create(defaultRole).Error; err != nil {
			return err
		}
	}

	return nil
}
