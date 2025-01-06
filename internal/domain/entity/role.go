package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code        string       `gorm:"size:64;uniqueIndex;not null" json:"code"`
	ParentID    *uint        `json:"parent_id"`                 // 父角色ID
	DataScope   string       `gorm:"size:32" json:"data_scope"` // 数据权限范围
	Status      int          `gorm:"default:1" json:"status"`
	Sort        int          `gorm:"default:0" json:"sort"` // 排序
	Description string       `gorm:"size:256" json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
}
