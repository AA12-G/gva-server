package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"type:varchar(64);uniqueIndex"`
	Password  string         `json:"-" gorm:"type:varchar(256)"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(128)"`
	Email     string         `json:"email" gorm:"type:varchar(128)"`
	Phone     string         `json:"phone" gorm:"type:varchar(32)"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(256)"`
	Status    int            `json:"status" gorm:"default:1"`
	RoleID    uint           `json:"role_id"`
	Role      *Role          `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
