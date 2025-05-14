package repository

import (
	"apac/internal/domain/entity"
	"apac/internal/domain/env"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthRepositoryItf interface {
	FindByEmail(string) (*entity.User, error)
	Create(*entity.User) error
	Update(string, *entity.User) error
	AddRefreshToken(userId uuid.UUID, token string) error
	GetUserRefreshTokens(userId uuid.UUID) ([]entity.RefreshToken, error)
	RemoveRefreshToken(token string) error
	FindByRefreshToken(token string) (*entity.User, error)
}

type AuthRepository struct {
	db  *gorm.DB
	env *env.Env
}

func NewAuthRepository(db *gorm.DB, env *env.Env) AuthRepositoryItf {
	return &AuthRepository{db, env}
}

func (r *AuthRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) Create(user *entity.User) error {
	user.PhotoURL = r.env.DefaultProfilePic
	return r.db.Create(user).Error
}

func (r *AuthRepository) Update(email string, user *entity.User) error {
	return r.db.Model(&entity.User{}).
		Where("email = ?", email).
		Updates(user).Error
}

func (r *AuthRepository) AddRefreshToken(userId uuid.UUID, token string) error {
	return r.db.Create(&entity.RefreshToken{
		Token:  token,
		UserID: userId,
	}).Error
}

func (r *AuthRepository) GetUserRefreshTokens(userId uuid.UUID) ([]entity.RefreshToken, error) {
	var refreshTokens []entity.RefreshToken
	err := r.db.Where("user_id = ?", userId).Find(&refreshTokens).Error
	return refreshTokens, err
}

func (r *AuthRepository) RemoveRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}

func (r *AuthRepository) FindByRefreshToken(token string) (*entity.User, error) {
	var refreshToken entity.RefreshToken
	err := r.db.Preload("User").Where("token = ?", token).First(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return refreshToken.User, nil
}
