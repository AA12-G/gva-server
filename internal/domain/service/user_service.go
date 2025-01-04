package service

import (
	"context"
	"errors"
	"gva/internal/domain/entity"
	"gva/internal/domain/repository"
	"gva/internal/infrastructure/cache"

	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
	cache    *cache.UserCache
}

func NewUserService(userRepo repository.UserRepository, db *gorm.DB, cache *cache.UserCache) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
		cache:    cache,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) error {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(ctx, username)
	if existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 获取默认角色
	var defaultRole entity.Role
	if err := s.db.Where("code = ?", "user").First(&defaultRole).Error; err != nil {
		return err
	}

	user := &entity.User{
		Username: username,
		Password: string(hashedPassword),
		Status:   1,
		RoleID:   defaultRole.ID, // 设置默认角色ID
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	return user, nil
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, nickname, email, phone, avatar string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Nickname = nickname
	user.Email = email
	user.Phone = phone
	user.Avatar = avatar

	return s.userRepo.Update(ctx, user)
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// UpdateAvatar 更新用户头像
func (s *UserService) UpdateAvatar(ctx context.Context, userID uint, avatarPath string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Avatar = avatarPath
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// 更新缓存
	return s.cache.DeleteUser(ctx, userID)
}

// GetUserByID 获取用户信息（使用缓存）
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	// 先从缓存获取
	user, err := s.cache.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// 缓存未命中，从数据库获取
	user, err = s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := s.cache.SetUser(ctx, user); err != nil {
		// 这里只记录日志，不返回错误
		log.Printf("缓存用户信息失败: %v", err)
	}

	return user, nil
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(ctx context.Context, page, pageSize int, keyword string, status *int) ([]*entity.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize, keyword, status)
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(ctx context.Context, userID uint, status int) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Status = status
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// 更新缓存
	return s.cache.DeleteUser(ctx, userID)
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	// 使用 GORM 的软删除功能
	if err := s.db.Delete(&entity.User{}, userID).Error; err != nil {
		return err
	}

	// 删除缓存
	return s.cache.DeleteUser(ctx, userID)
}
