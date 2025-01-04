package service

import (
	"context"
	"gva/internal/domain/entity"
	"gva/internal/domain/repository"
)

type OperationLogService struct {
	logRepo repository.OperationLogRepository
}

func NewOperationLogService(logRepo repository.OperationLogRepository) *OperationLogService {
	return &OperationLogService{logRepo: logRepo}
}

func (s *OperationLogService) Create(ctx context.Context, log *entity.OperationLog) error {
	return s.logRepo.Create(ctx, log)
}

func (s *OperationLogService) List(ctx context.Context, page, size int) ([]*entity.OperationLog, int64, error) {
	return s.logRepo.List(ctx, page, size)
}
