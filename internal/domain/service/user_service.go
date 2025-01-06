package service

import (
	"context"
	"errors"
	"gva/internal/domain/entity"
	"gva/internal/domain/repository"
	"gva/internal/infrastructure/cache"
	"gva/internal/pkg/jwt"
	"gva/internal/pkg/utils"

	"log"

	"encoding/csv"
	"fmt"
	"io"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 定义登录相关的错误
var (
	ErrUserNotFound    = errors.New("用户不存在")
	ErrUserDisabled    = errors.New("用户已被禁用")
	ErrIncorrectPass   = errors.New("密码错误")
	ErrUserNotVerified = errors.New("用户未通过审核")
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
		// 如果找不到默认角色，创建一个
		defaultRole = entity.Role{
			Name: "普通用户",
			Code: "user",
		}
		if err := s.db.Create(&defaultRole).Error; err != nil {
			return err
		}
	}

	user := &entity.User{
		Username: username,
		Password: string(hashedPassword),
		Status:   1,
		RoleID:   defaultRole.ID,
	}

	return s.userRepo.Create(ctx, user)
}

// Login 用户登录
func (s *UserService) Login(ctx context.Context, username, password string) (*entity.User, string, error) {
	// 缓存未命中，从数据库查询
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrUserNotFound
		}
		return nil, "", fmt.Errorf("查询用户失败: %v", err)
	}

	// 将用户信息存入缓存
	if err := s.cache.SetUser(ctx, user); err != nil {
		log.Printf("缓存用户信息失败: %v", err)
	}

	// 检查用户状态
	switch user.Status {
	case 0:
		return nil, "", ErrUserDisabled
	case 2:
		return nil, "", ErrUserNotVerified
	}

	fmt.Println("user.Password", user.Password)

	// 验证密码
	if !utils.CheckPassword(password, user.Password) {
		return nil, "", ErrIncorrectPass
	}

	// 生成token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("生成token失败: %v", err)
	}

	// 预加载角色信息
	if err := s.db.Preload("Role").First(user, user.ID).Error; err != nil {
		return nil, "", fmt.Errorf("加载用户角色失败: %v", err)
	}

	return user, token, nil
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
	return s.cache.DeleteUserByID(ctx, userID)
}

// GetUserByID 获取用户信息（使用缓存）
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	// 先从缓存获取
	user, err := s.cache.GetUserByID(ctx, id)
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
	return s.cache.DeleteUserByID(ctx, userID)
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	// 先获取用户信息，用于删除缓存
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 使用 GORM 的软删除功能
	if err := s.db.Delete(&entity.User{}, userID).Error; err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}

	// 删除缓存
	if err := s.cache.DeleteUserByID(ctx, userID); err != nil {
		log.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// ExportUsers 导出所有用户
func (s *UserService) ExportUsers(ctx context.Context) ([]*entity.User, error) {
	return s.userRepo.FindAll(ctx)
}

// ImportUsers 导入用户
func (s *UserService) ImportUsers(ctx context.Context, reader io.Reader) ([]*entity.User, error) {
	// 获取默认角色
	var defaultRole entity.Role
	if err := s.db.Where("code = ?", "user").First(&defaultRole).Error; err != nil {
		// 如果找不到默认角色，创建一个
		defaultRole = entity.Role{
			Name: "普通用户",
			Code: "user",
		}
		if err := s.db.Create(&defaultRole).Error; err != nil {
			return nil, err
		}
	}

	r := csv.NewReader(reader)
	r.FieldsPerRecord = -1 // 允许字段数量不固定

	// 跳过表头
	if _, err := r.Read(); err != nil {
		return nil, err
	}

	var users []*entity.User
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 4 { // 至少需要用户名、昵称、邮箱、手机号
			continue
		}

		// 生成随机密码
		password := fmt.Sprintf("%08d", rand.Intn(100000000))
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entity.User{
			Username: record[0],
			Nickname: record[1],
			Email:    record[2],
			Phone:    record[3],
			Password: string(hashedPassword),
			Status:   1,              // 默认正常状态
			RoleID:   defaultRole.ID, // 设置默认角色ID
		}

		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserWithRole 获取用户信息（包含角色和权限）
func (s *UserService) GetUserWithRole(ctx context.Context, userID uint) (*entity.User, error) {
	var user entity.User
	err := s.db.Preload("Role").
		Preload("Role.Permissions").
		First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}

	return &user, nil
}
