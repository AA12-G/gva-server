package entity

import "gorm.io/gorm"

// Dict 数据字典
type Dict struct {
	gorm.Model
	Name        string     `gorm:"size:64;uniqueIndex;not null" json:"name"` // 字典名称
	Code        string     `gorm:"size:64;uniqueIndex;not null" json:"code"` // 字典编码
	Description string     `gorm:"size:256" json:"description"`              // 描述
	Status      int        `gorm:"default:1" json:"status"`                  // 状态
	Items       []DictItem `json:"items"`                                    // 字典项
}

// DictItem 字典项
type DictItem struct {
	gorm.Model
	DictID      uint   `json:"dict_id"`                        // 所属字典ID
	Label       string `gorm:"size:128;not null" json:"label"` // 标签
	Value       string `gorm:"size:128;not null" json:"value"` // 值
	Sort        int    `gorm:"default:0" json:"sort"`          // 排序
	Status      int    `gorm:"default:1" json:"status"`        // 状态
	Description string `gorm:"size:256" json:"description"`    // 描述
}
