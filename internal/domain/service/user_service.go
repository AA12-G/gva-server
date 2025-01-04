package service

import (
	"context"
	"errors"
	"gva/internal/domain/entity"
	"gva/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewUserService(userRepo repository.UserRepository, db *gorm.DB) *UserService {
	return &UserService{
		userRepo: userRepo,
		db:       db,
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
