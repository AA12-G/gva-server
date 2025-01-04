package entity

import (
	"time"
)

// OperationLog 操作日志
type OperationLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user" gorm:"foreignKey:UserID"`
	IP        string    `json:"ip"`
	Method    string    `json:"method"`
	Path      string    `json:"path"`
	Status    int       `json:"status"`
	Latency   int64     `json:"latency"` // 请求耗时（毫秒）
	UserAgent string    `json:"user_agent"`
	Request   string    `json:"request"`  // 请求参数
	Response  string    `json:"response"` // 响应内容
	CreatedAt time.Time `json:"created_at"`
}
