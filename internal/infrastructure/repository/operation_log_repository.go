package repository

import (
	"context"
	"gva/internal/domain/entity"

	"gorm.io/gorm"
)

type operationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) *operationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(ctx context.Context, log *entity.OperationLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *operationLogRepository) List(ctx context.Context, page, size int) ([]*entity.OperationLog, int64, error) {
	var logs []*entity.OperationLog
	var total int64

	db := r.db.WithContext(ctx)

	// 统计总数
	if err := db.Model(&entity.OperationLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Preload("User"). // 预加载用户信息
					Order("id DESC").
					Offset((page - 1) * size).
					Limit(size).
					Find(&logs).Error

	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
