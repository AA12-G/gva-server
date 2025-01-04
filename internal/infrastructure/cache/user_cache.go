package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gva/internal/domain/entity"
	"time"

	"github.com/go-redis/redis/v8"
)

type UserCache struct {
	rdb *redis.Client
}

func NewUserCache(rdb *redis.Client) *UserCache {
	return &UserCache{rdb: rdb}
}

// SetUser 缓存用户信息
func (c *UserCache) SetUser(ctx context.Context, user *entity.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetUser 获取缓存的用户信息
func (c *UserCache) GetUser(ctx context.Context, id uint) (*entity.User, error) {
	key := fmt.Sprintf("user:%d", id)
	data, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser 删除用户缓存
func (c *UserCache) DeleteUser(ctx context.Context, id uint) error {
	key := fmt.Sprintf("user:%d", id)
	return c.rdb.Del(ctx, key).Err()
}
