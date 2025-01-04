package database

import (
	"fmt"
	"time"

	"gva/internal/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLConnection(config config.MySQLConfig) (*gorm.DB, error) {
	// 打印配置
	fmt.Printf("MySQL config: %+v\n", config)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
