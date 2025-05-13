package repository

import (
	"apac/internal/domain/entity"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryItf interface {
	FindById(userId uuid.UUID) (*entity.User, error)
	UpdateUser(userId uuid.UUID, user *entity.User) error
	AddPreference(userId uuid.UUID, preference string) error
	RemovePreference(userId uuid.UUID, preference string) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryItf {
	return &UserRepository{db}
}

func (r *UserRepository) FindById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", userId).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(userId uuid.UUID, user *entity.User) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userId).Updates(user).Error
}

func (r *UserRepository) AddPreference(userId uuid.UUID, preference string) error {
	return r.db.Create(&entity.Preference{
		UserID: userId,
		Name:   preference,
	}).Error
}

func (r *UserRepository) RemovePreference(userId uuid.UUID, preference string) error {
	return r.db.Delete(&entity.Preference{}, "user_id = ? AND name = ?", userId, preference).Error
}
