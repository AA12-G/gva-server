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

// GetUserByID 通过ID获取用户
func (c *UserCache) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	key := fmt.Sprintf("user:id:%d", id)
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

// GetUserByUsername 通过用户名获取用户
func (c *UserCache) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	key := fmt.Sprintf("user:username:%s", username)
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

// SetUser 缓存用户信息
func (c *UserCache) SetUser(ctx context.Context, user *entity.User) error {

	// 设置两个缓存键，一个用ID索引，一个用用户名索引
	idKey := fmt.Sprintf("user:id:%d", user.ID)
	usernameKey := fmt.Sprintf("user:username:%s", user.Username)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	pipe := c.rdb.Pipeline()
	pipe.Set(ctx, idKey, data, time.Hour)
	pipe.Set(ctx, usernameKey, data, time.Hour)
	_, err = pipe.Exec(ctx)
	return err
}

// DeleteUser 删除用户缓存
func (c *UserCache) DeleteUser(ctx context.Context, user *entity.User) error {
	idKey := fmt.Sprintf("user:id:%d", user.ID)
	usernameKey := fmt.Sprintf("user:username:%s", user.Username)

	pipe := c.rdb.Pipeline()
	pipe.Del(ctx, idKey)
	pipe.Del(ctx, usernameKey)
	_, err := pipe.Exec(ctx)
	return err
}

// DeleteUserByID 通过ID删除用户缓存
func (c *UserCache) DeleteUserByID(ctx context.Context, id uint) error {
	// 先获取用户信息，以便删除username索引
	user, err := c.GetUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		// 用户不存在，直接删除id索引
		idKey := fmt.Sprintf("user:id:%d", id)
		return c.rdb.Del(ctx, idKey).Err()
	}

	// 删除所有相关缓存
	return c.DeleteUser(ctx, user)
}
