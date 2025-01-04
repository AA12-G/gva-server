package repository

import (
	"context"
	"gva/internal/domain/entity"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, page, size int, keyword string, status *int) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	db := r.db.WithContext(ctx)

	// 构建查询条件
	if keyword != "" {
		db = db.Where("username = ? OR nickname LIKE ? OR email LIKE ?",
			keyword, "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		db = db.Where("status = ?", *status)
	}

	// 统计总数
	if err := db.Model(&entity.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := db.Preload("Role"). // 预加载角色信息
					Offset((page - 1) * size).
					Limit(size).
					Order("id DESC").
					Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
