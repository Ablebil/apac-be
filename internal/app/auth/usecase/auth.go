package usecase

import (
	"apac/internal/app/user/repository"
	"apac/internal/domain/dto"
	"apac/internal/domain/entity"
	"apac/internal/domain/env"
	"apac/internal/infra/code"
	"apac/internal/infra/oauth"
	"apac/internal/infra/response"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthUsecaseItf interface {
	GoogleLogin() (string, *response.Err)
	GoogleCallback(data dto.GoogleCallbackRequest) (string, *response.Err)
}

type AuthUsecase struct {
	UserRepository repository.UserRepositoryItf
	code           code.CodeItf
	db             *gorm.DB
	env            *env.Env
	OAuth          oauth.OAuthItf
}

func NewAuthUsecase(userRepository repository.UserRepositoryItf, code code.CodeItf, db *gorm.DB, env *env.Env, oauth oauth.OAuthItf) AuthUsecaseItf {
	return &AuthUsecase{
		UserRepository: userRepository,
		code:           code,
		db:             db,
		env:            env,
		OAuth:          oauth,
	}
}

func (a *AuthUsecase) GoogleLogin() (string, *response.Err) {
	state, err := a.code.GenerateToken()
	if err != nil {
		return "", response.ErrInternalServer()
	}

	url, err := a.OAuth.GenerateAuthLink(state)
	if err != nil {
		return "", response.ErrInternalServer()
	}

	return url, nil
}

func (a *AuthUsecase) GoogleCallback(data dto.GoogleCallbackRequest) (string, *response.Err) {
	tx := a.db.Begin()

	defer func() {
		tx.Rollback()
	}()

	if data.Error != "" {
		return "", response.ErrForbidden(data.Error)
	}

	token, err := a.OAuth.ExchangeToken(data.Code)
	if err != nil {
		return "", response.ErrInternalServer()
	}

	profile, err := a.OAuth.GetProfile(token)
	if err != nil {
		return "", response.ErrInternalServer()
	}

	id, _ := uuid.NewV7()

	user := &entity.User{
		ID:         id,
		Name:       profile.Name,
		Username:   profile.Username,
		Email:      profile.Email,
		IsVerified: profile.IsVerified,
	}
}
