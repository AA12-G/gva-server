package database

import (
	"gva/internal/pkg/config"
	"log"

	"gorm.io/gorm"
)

// InitTestDB 初始化测试数据库
func InitTestDB() *gorm.DB {
	conf := &config.MySQLConfig{
		Host:         "localhost",
		Port:         3306,
		Username:     "root",
		Password:     "root",
		Database:     "gva_test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
	}

	db, err := InitDB(conf) // 修改这里，使用新的 InitDB 函数
	if err != nil {
		log.Fatalf("初始化测试数据库失败: %v", err)
	}

	return db
}

// CleanTestDB 清理测试数据库
func CleanTestDB(db *gorm.DB) {
	// 清理所有表数据
	tables := []string{"users", "roles", "permissions", "role_permissions", "operation_logs"}
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}
}
