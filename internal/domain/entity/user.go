package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"size:256;not null" json:"-"`
	Nickname  string    `gorm:"size:128" json:"nickname"`
	Email     string    `gorm:"size:128" json:"email"`
	Phone     string    `gorm:"size:32" json:"phone"`
	Avatar    string    `gorm:"size:256" json:"avatar"`
	Status    int       `gorm:"default:1" json:"status"` // 1: 正常, 0: 禁用
	LastLogin time.Time `json:"last_login"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
}
