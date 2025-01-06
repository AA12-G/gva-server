package database

import (
	"fmt"
	"gva/internal/domain/entity"
	"gva/internal/pkg/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(conf *config.MySQLConfig) (*gorm.DB, error) {
	log.Printf("正在连接数据库: %s@%s:%d/%s", conf.Username, conf.Host, conf.Port, conf.Database)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库实例失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 执行自动迁移
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("自动迁移失败: %v", err)
	}

	// 检查是否需要初始化数据
	var count int64
	db.Model(&entity.Role{}).Count(&count)
	if count == 0 {
		log.Println("开始初始化基础数据...")
		if err := InitData(db); err != nil {
			return nil, fmt.Errorf("初始化数据失败: %v", err)
		}
		log.Println("基础数据初始化完成")
	}

	return db, nil
}

// AutoMigrate 自动迁移数据库结构
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.OperationLog{},
		// 添加其他需要迁移的实体
	)
}
