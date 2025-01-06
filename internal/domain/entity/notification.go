package entity

import "gorm.io/gorm"

// Notification 通知
type Notification struct {
	gorm.Model
	Title     string `gorm:"size:128;not null" json:"title"`                 // 标题
	Content   string `gorm:"type:text" json:"content"`                       // 内容
	Type      string `gorm:"size:32" json:"type"`                            // 通知类型
	Status    int    `gorm:"default:0" json:"status"`                        // 状态
	SenderID  uint   `json:"sender_id"`                                      // 发送者ID
	Sender    User   `json:"sender"`                                         // 发送者
	Receivers []User `gorm:"many2many:user_notifications;" json:"receivers"` // 接收者
}

// UserNotification 用户通知关联
type UserNotification struct {
	gorm.Model
	UserID         uint  `json:"user_id"`                      // 用户ID
	NotificationID uint  `json:"notification_id"`              // 通知ID
	ReadStatus     int   `gorm:"default:0" json:"read_status"` // 阅读状态
	ReadTime       *Time `json:"read_time"`                    // 阅读时间
}
