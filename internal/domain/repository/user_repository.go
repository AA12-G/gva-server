package repository

import (
	"context"
	"gva/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	List(ctx context.Context, page, size int, keyword string, status *int) ([]*entity.User, int64, error)
}
