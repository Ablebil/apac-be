package usecase

import (
	trepo "apac/internal/app/trip/repository"
	urepo "apac/internal/app/user/repository"
	"apac/internal/domain/dto"
	"apac/internal/domain/entity"
	"apac/internal/domain/env"
	"apac/internal/infra/gemini"
	res "apac/internal/infra/response"
	"encoding/json"

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
	var user *entity.User
	if userId != uuid.Nil {
		foundUser, err := uc.userRepository.FindById(userId)
		user = foundUser

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

	content, err := json.Marshal(response)
	if err != nil {
		return nil, res.ErrInternalServer("Unable to parse JSON response into string")
	}

	trip := &entity.Trip{
		UserID:  userId,
		Content: string(content),
		User:    user,
	}

	if err := uc.tripRepository.Create(trip); err != nil {
		return nil, res.ErrInternalServer("Cannot add trip to history: " + err.Error())
	}

	return response, nil
}
