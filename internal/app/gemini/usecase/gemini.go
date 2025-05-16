package usecase

import (
	trepo "apac/internal/app/trip/repository"
	urepo "apac/internal/app/user/repository"
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
	userRepository urepo.UserRepositoryItf
	tripRepository trepo.TripRepositoryItf
}

func NewGeminiUsecase(
	env *env.Env,
	gemini gemini.GeminiItf,
	userRepository urepo.UserRepositoryItf,
	tripRepository trepo.TripRepositoryItf,
) GeminiUsecaseItf {
	return &GeminiUsecase{
		env:            env,
		gemini:         gemini,
		userRepository: userRepository,
		tripRepository: tripRepository,
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

	if err := uc.tripRepository.Create(response); err != nil {
		return nil, res.ErrInternalServer("Cannot add trip to history")
	}

	return response, nil
}
