package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string       `json:"name" gorm:"size:64;uniqueIndex;not null"`
	Code        string       `json:"code" gorm:"size:64;uniqueIndex;not null"`
	ParentID    *uint        `json:"parent_id"`  // 父角色ID
	DataScope   string       `json:"data_scope"` // 数据权限范围
	Status      int          `json:"status" gorm:"default:1"`
	Sort        int          `json:"sort" gorm:"default:0"`
	Description string       `json:"description" gorm:"size:256"`
	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"` // 角色拥有的权限
}
