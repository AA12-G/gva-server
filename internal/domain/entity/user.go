package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string         `json:"username" gorm:"size:64;uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"size:128;not null"`
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
