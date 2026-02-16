package repo

import (
	"errors"
	"wyy/internal/domain"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// GetByID 根据 ID 查询用户
func (r *UserRepo) GetByID(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // 或返回自定义错误
	}
	return &user, err
}

// GetByEmail 根据邮箱查询用户
func (r *UserRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Create 创建新用户
func (r *UserRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除或硬删除）
func (r *UserRepo) Delete(id int64) error {
	return r.db.Delete(&domain.User{}, id).Error
}
