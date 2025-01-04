package repository

import (
	"context"
	"gva/internal/domain/entity"
)

// OperationLogRepository 操作日志仓储接口
type OperationLogRepository interface {
	Create(ctx context.Context, log *entity.OperationLog) error
	List(ctx context.Context, page, size int) ([]*entity.OperationLog, int64, error)
}
