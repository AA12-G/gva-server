package database

import (
	"gva/internal/pkg/config"
	"log"

	"gorm.io/gorm"
)

// InitTestDB 初始化测试数据库
func InitTestDB() *gorm.DB {
	cfg := &config.MySQLConfig{
		Host:         "localhost",
		Port:         3306,
		Username:     "root",
		Password:     "root",
		Database:     "gva_test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}

	db, err := NewMySQLDB(cfg)
	if err != nil {
		log.Fatalf("连接测试数据库失败: %v", err)
	}

	// 自动迁移数据库
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("测试数据库迁移失败: %v", err)
	}

	return db
}

// AutoMigrate 自动迁移数据库
func AutoMigrate(db *gorm.DB) error {
	// 在这里添加需要迁移的模型
	return nil
}

// CleanTestDB 清理测试数据库
func CleanTestDB(db *gorm.DB) {
	// 清理所有表数据
	tables := []string{"users", "roles", "permissions", "role_permissions", "operation_logs"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}
}
