package usecase

import (
	"apac/internal/domain/dto"
	"apac/internal/domain/env"
	"apac/internal/infra/gemini"
	res "apac/internal/infra/response"
)

type GeminiUsecaseItf interface {
	Prompt(*dto.GeminiRequest) (map[string]interface{}, *res.Err)
}

type GeminiUsecase struct {
	env    *env.Env
	gemini gemini.GeminiItf
}

func NewGeminiUsecase(env *env.Env, gemini gemini.GeminiItf) GeminiUsecaseItf {
	return &GeminiUsecase{
		env:    env,
		gemini: gemini,
	}
}

func (uc *GeminiUsecase) Prompt(payload *dto.GeminiRequest) (map[string]interface{}, *res.Err) {
	response, err := uc.gemini.Prompt(payload.Text)
	if err != nil {
		return nil, res.ErrInternalServer("AI prompting failed")
	}

	return response, nil
}
