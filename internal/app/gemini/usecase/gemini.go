package usecase

import (
	"apac/internal/app/user/repository"
	"apac/internal/domain/dto"
	"apac/internal/domain/env"
	"apac/internal/infra/gemini"
	res "apac/internal/infra/response"

	"github.com/google/uuid"
)

type GeminiUsecaseItf interface {
	Prompt(*dto.GeminiRequest, uuid.UUID) (map[string]interface{}, *res.Err)
}

type GeminiUsecase struct {
	env            *env.Env
	gemini         gemini.GeminiItf
	userRepository repository.UserRepositoryItf
}

func NewGeminiUsecase(env *env.Env, gemini gemini.GeminiItf, userRepository repository.UserRepositoryItf) GeminiUsecaseItf {
	return &GeminiUsecase{
		env:            env,
		gemini:         gemini,
		userRepository: userRepository,
	}
}

func (uc *GeminiUsecase) Prompt(payload *dto.GeminiRequest, userId uuid.UUID) (map[string]interface{}, *res.Err) {
	var preferences []string
	if userId != uuid.Nil {
		user, err := uc.userRepository.FindById(userId)

		if err != nil {
			return nil, res.ErrInternalServer("Failed to find user")
		}

		if user == nil {
			return nil, res.ErrNotFound("User not found")
		}

		preferences = user.ParseDTOGet().Preferences
	}

	response, err := uc.gemini.Prompt(preferences, payload.Text)
	if err != nil {
		return nil, res.ErrInternalServer("AI prompting failed: " + err.Error())
	}

	return response, nil
}
