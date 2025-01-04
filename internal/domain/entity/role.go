package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code        string       `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description string       `gorm:"size:256" json:"description"`
	Status      int          `gorm:"default:1" json:"status"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
