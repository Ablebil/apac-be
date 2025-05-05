package repository

import (
	"apac/internal/domain/entity"

	"gorm.io/gorm"
)

type UserRepositoryItf interface {
	Create(tx *gorm.DB, user *entity.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryItf {
	return &UserRepository{db}
}

func (r *UserRepository) Create(tx *gorm.DB, user *entity.User) error {
	return r.db.Debug().Create(user).Error
}
