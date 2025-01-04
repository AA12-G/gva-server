package database

import (
	"gva/internal/domain/entity"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库
func AutoMigrate(db *gorm.DB) error {
	log.Println("开始迁移数据库...")

	// 创建表
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
		&entity.OperationLog{},
	)
	if err != nil {
		return err
	}

	// 创建默认角色
	var count int64
	db.Model(&entity.Role{}).Count(&count)
	if count == 0 {
		roles := []entity.Role{
			{
				Name: "管理员",
				Code: "admin",
			},
			{
				Name: "普通用户",
				Code: "user",
			},
		}
		if err := db.Create(&roles).Error; err != nil {
			return err
		}
	}

	return nil
}
