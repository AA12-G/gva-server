package entity

import "gorm.io/gorm"

// Department 部门
type Department struct {
	gorm.Model
	Name        string `gorm:"size:64;not null" json:"name"`    // 部门名称
	Code        string `gorm:"size:64;uniqueIndex" json:"code"` // 部门编码
	ParentID    *uint  `json:"parent_id"`                       // 父部门ID
	Path        string `gorm:"size:255" json:"path"`            // 部门路径
	Leader      string `gorm:"size:64" json:"leader"`           // 部门负责人
	Phone       string `gorm:"size:32" json:"phone"`            // 联系电话
	Email       string `gorm:"size:128" json:"email"`           // 邮箱
	Sort        int    `gorm:"default:0" json:"sort"`           // 排序
	Status      int    `gorm:"default:1" json:"status"`         // 状态
	Description string `gorm:"size:256" json:"description"`     // 描述
	Users       []User `json:"users"`                           // 部门用户
}
