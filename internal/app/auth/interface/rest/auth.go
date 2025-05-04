package rest

import (
	"apac/internal/app/auth/usecase"
	"apac/internal/domain/dto"
	"apac/internal/infra/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Validator   *validator.Validate
	AuthUsecase usecase.AuthUsecaseItf
}

func NewAuthHandler(router fiber.Router, usecase usecase.AuthUsecaseItf, validator *validator.Validate) {
	authHandler := AuthHandler{
		Validator:   validator,
		AuthUsecase: usecase,
	}

	routerGroup := router.Group("/auth")
	routerGroup.Get("/google", authHandler.GoogleLogin)
	routerGroup.Post("/google", authHandler.GoogleCallback)
}

func (h AuthHandler) GoogleLogin(ctx *fiber.Ctx) error {
	url, err := h.AuthUsecase.GoogleLogin()
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.SuccessResponse(ctx, response.OAuthLoginSuccess, fiber.Map{
		"url": url,
	})
}

func (h AuthHandler) GoogleCallback(ctx *fiber.Ctx) error {
	data := new(dto.GoogleCallbackRequest)
	if err := ctx.BodyParser(data); err != nil {
		return response.BadRequest(ctx)
	}

	if err := h.Validator.Struct(data); err != nil {
		return response.ValidationError(ctx, err)
	}

	token, err := h.AuthUsecase.GoogleCallback(*data)
	if err != nil {
		return response.Error(ctx, err)
	}

	return response.SuccessResponse(ctx, response.LoginSuccess, fiber.Map{
		"token": token,
	})
}
