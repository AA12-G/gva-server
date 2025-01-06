package entity

import "gorm.io/gorm"

// Task 定时任务
type Task struct {
	gorm.Model
	Name        string `gorm:"size:64;not null" json:"name"`     // 任务名称
	Cron        string `gorm:"size:64;not null" json:"cron"`     // cron表达式
	Command     string `gorm:"size:512;not null" json:"command"` // 执行命令
	Status      int    `gorm:"default:1" json:"status"`          // 状态
	Description string `gorm:"size:256" json:"description"`      // 描述
	LastRunTime *Time  `json:"last_run_time"`                    // 上次执行时间
	NextRunTime *Time  `json:"next_run_time"`                    // 下次执行时间
}

// TaskLog 任务执行日志
type TaskLog struct {
	gorm.Model
	TaskID    uint   `json:"task_id"`                 // 任务ID
	Status    int    `json:"status"`                  // 执行状态
	Result    string `gorm:"type:text" json:"result"` // 执行结果
	StartTime Time   `json:"start_time"`              // 开始时间
	EndTime   Time   `json:"end_time"`                // 结束时间
}
