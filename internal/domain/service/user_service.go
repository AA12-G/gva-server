package service

import (
	"context"
	"errors"
	"gva/internal/domain/cache"
	"gva/internal/domain/entity"
	"gva/internal/domain/repository"
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
	ErrUserNotFound         = errors.New("用户不存在")
	ErrUserDisabled         = errors.New("用户已被禁用")
	ErrIncorrectPass        = errors.New("密码错误")
	ErrUserNotVerified      = errors.New("用户未通过审核")
	ErrUserFrozen           = errors.New("您的账号已被冻结，请联系客服")
	ErrUserPhoneNotVerified = errors.New("您的账号未绑定手机号或手机号码格式不正确，请联系客服")
)

type UserService struct {
	userRepo repository.UserRepository
	db       *gorm.DB
	cache    cache.UserCache
}

func NewUserService(userRepo repository.UserRepository, db *gorm.DB, cache cache.UserCache) *UserService {
	// 如果缓存为空，创建一个空实现
	if cache == nil {
		cache = &EmptyCache{}
	}
	return &UserService{
		userRepo: userRepo,
		db:       db,
		cache:    cache,
	}
}

// EmptyCache 空缓存实现
type EmptyCache struct{}

// 实现 cache.UserCache 接口
func (c *EmptyCache) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return nil, nil
}

func (c *EmptyCache) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return nil, nil
}

func (c *EmptyCache) SetUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (c *EmptyCache) DeleteUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (c *EmptyCache) DeleteUserByID(ctx context.Context, id uint) error {
	return nil
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

	if user.Status == 2 {
		// 您的账号已被冻结，请联系客服
		return nil, "", ErrUserFrozen
	}

	if user.Phone == "" && len(user.Phone) < 11 && len(user.Phone) > 11 {
		// 您的账号未绑定手机号或手机号码格式不正确，请联系客服
		return nil, "", ErrUserPhoneNotVerified
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

	fmt.Println("222user.Password", user.Password)

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
func (s *UserService) UpdateProfile(ctx context.Context, userID uint, username, nickname, email, phone string, roleID uint) error {
	// 获取用户信息
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}

	// 如果用户名发生变化，检查新用户名是否已存在
	if username != user.Username {
		existingUser, err := s.userRepo.FindByUsername(ctx, username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("检查用户名失败: %v", err)
		}
		if existingUser != nil {
			return fmt.Errorf("用户名 %s 已被使用", username)
		}
	}

	// 如果要更新角色，先检查角色是否存在
	if roleID > 0 {
		var role entity.Role
		if err := s.db.First(&role, roleID).Error; err != nil {
			return fmt.Errorf("角色不存在或已被删除")
		}
		user.RoleID = roleID
	}

	// 更新用户信息
	user.Username = username
	user.Nickname = nickname
	user.Email = email
	user.Phone = phone

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("更新用户信息失败: %v", err)
	}

	// 删除缓存
	if err := s.cache.DeleteUserByID(ctx, userID); err != nil {
		log.Printf("删除用户缓存失败: %v", err)
	}

	return nil
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
	log.Printf("开始更新用户状态 - UserID: %d, Status: %d", userID, status)

	var user entity.User
	if err := s.db.First(&user, userID).Error; err != nil {
		log.Printf("查询用户失败 - UserID: %d, Error: %v", userID, err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %v", err)
	}

	log.Printf("找到用户 - UserID: %d, Username: %s", user.ID, user.Username)

	if err := s.db.Model(&user).Update("status", status).Error; err != nil {
		log.Printf("更新状态失败 - UserID: %d, Error: %v", userID, err)
		return fmt.Errorf("更新状态失败: %v", err)
	}

	log.Printf("更新状态成功 - UserID: %d, NewStatus: %d", userID, status)
	return nil
}

// DeleteUser 删除用户（硬删除）
func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	// 先检查用户是否存在
	if err := s.db.First(&entity.User{}, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 使用 Unscoped 来执行硬删除
	if err := s.db.Unscoped().Delete(&entity.User{}, userID).Error; err != nil {
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

// UpdateUserRole 更新用户角色
func (s *UserService) UpdateUserRole(ctx context.Context, userID uint, roleID uint) error {
	// 检查用户是否存在
	var user entity.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %v", err)
	}

	// 检查角色是否存在
	var role entity.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("角色不存在")
		}
		return fmt.Errorf("查询角色失败: %v", err)
	}

	// 更新用户角色
	if err := s.db.Model(&user).Update("role_id", roleID).Error; err != nil {
		return fmt.Errorf("更新用户角色失败: %v", err)
	}

	// 删除用户缓存
	if err := s.cache.DeleteUserByID(ctx, userID); err != nil {
		log.Printf("删除用户缓存失败: %v", err)
	}

	return nil
}

// RestoreUser 恢复已删除的用户
func (s *UserService) RestoreUser(ctx context.Context, userID uint) error {
	// 使用 Unscoped 查找被删除的用户
	var user entity.User
	if err := s.db.Unscoped().First(&user, userID).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 检查用户是否已被删除
	if user.DeletedAt.Time.IsZero() {
		return fmt.Errorf("用户未被删除")
	}

	// 恢复用户
	if err := s.db.Unscoped().Model(&user).Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("恢复用户失败: %v", err)
	}

	// 更新缓存
	if err := s.cache.SetUser(ctx, &user); err != nil {
		log.Printf("更新用户缓存失败: %v", err)
	}

	return nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, username, password, nickname, email, phone string, roleID uint) error {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查用户名失败: %v", err)
	}
	if existingUser != nil {
		return fmt.Errorf("用户名 %s 已存在", username)
	}

	// 检查角色是否存在
	if roleID > 0 {
		var role entity.Role
		if err := s.db.First(&role, roleID).Error; err != nil {
			return fmt.Errorf("角色不存在或已被删除")
		}
	}

	// 如果没有提供密码，使用默认密码
	if password == "" {
		password = "123456"
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("生成密码失败: %v", err)
	}

	// 创建用户
	user := &entity.User{
		Username: username,
		Password: hashedPassword,
		Nickname: nickname,
		Email:    email,
		Phone:    phone,
		RoleID:   roleID,
		Status:   1, // 默认启用
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}

	return nil
}
