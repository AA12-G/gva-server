package database

import (
	"gva/internal/infrastructure/config"

	"gorm.io/gorm"
)

var testDB *gorm.DB

func GetTestDB() (*gorm.DB, error) {
	if testDB != nil {
		return testDB, nil
	}

	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := NewMySQLConnection(cfg.MySQL)
	if err != nil {
		return nil, err
	}

	testDB = db
	return db, nil
}
