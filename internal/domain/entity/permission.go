package entity

import "gorm.io/gorm"

const (
	MenuPermission   = "menu"   // 菜单权限
	ButtonPermission = "button" // 按钮权限
	DataPermission   = "data"   // 数据权限
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code        string `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Type        string `gorm:"size:32;not null" json:"type"` // 权限类型
	ParentID    *uint  `json:"parent_id"`                    // 父权限ID
	Path        string `gorm:"size:128" json:"path"`         // 路由路径
	Component   string `gorm:"size:128" json:"component"`    // 前端组件
	Redirect    string `gorm:"size:128" json:"redirect"`     // 重定向地址
	Icon        string `gorm:"size:64" json:"icon"`          // 图标
	Sort        int    `gorm:"default:0" json:"sort"`        // 排序
	Hidden      bool   `gorm:"default:false" json:"hidden"`  // 是否隐藏
	Status      int    `gorm:"default:1" json:"status"`      // 状态
	Description string `gorm:"size:256" json:"description"`  // 描述
}
