package cache

import (
	"context"
	"gva/internal/domain/entity"
)

// UserCache 用户缓存接口
type UserCache interface {
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	SetUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, user *entity.User) error
	DeleteUserByID(ctx context.Context, id uint) error
}
