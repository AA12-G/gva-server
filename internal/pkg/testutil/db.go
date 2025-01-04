package testutil

import (
	infraconfig "gva/internal/infrastructure/config"
	"gva/internal/infrastructure/database"
	"gva/internal/pkg/config"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func GetTestDB() (*gorm.DB, error) {
	if testDB != nil {
		return testDB, nil
	}

	cfg, err := infraconfig.Load()
	if err != nil {
		return nil, err
	}

	mysqlConfig := config.MySQLConfig{
		Host:         cfg.MySQL.Host,
		Port:         cfg.MySQL.Port,
		Username:     cfg.MySQL.Username,
		Password:     cfg.MySQL.Password,
		Database:     cfg.MySQL.Database,
		MaxIdleConns: cfg.MySQL.MaxIdleConns,
		MaxOpenConns: cfg.MySQL.MaxOpenConns,
	}

	db, err := database.NewMySQLConnection(mysqlConfig)
	if err != nil {
		return nil, err
	}

	testDB = db
	return db, nil
}
