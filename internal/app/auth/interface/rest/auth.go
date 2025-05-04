package rest

import (
	"apac/internal/app/auth/usecase"
	"apac/internal/domain/dto"
	res "apac/internal/infra/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	Validator   *validator.Validate
	AuthUsecase usecase.AuthUsecaseItf
}

func NewAuthHandler(routerGroup fiber.Router, authUsecase usecase.AuthUsecaseItf, validator *validator.Validate) {
	AuthHandler := AuthHandler{
		Validator:   validator,
		AuthUsecase: authUsecase,
	}

	routerGroup = routerGroup.Group("/auth")
	routerGroup.Post("/register", AuthHandler.Register)
	routerGroup.Post("/verify-otp", AuthHandler.VerifyOTP)
	routerGroup.Post("/login", AuthHandler.Login)
}

func (h AuthHandler) Register(ctx *fiber.Ctx) error {
	payload := new(dto.RegisterRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.Register(payload); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "OTP sent to email", nil)
}

func (h AuthHandler) VerifyOTP(ctx *fiber.Ctx) error {
	payload := new(dto.VerifyOTPRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	accessToken, refreshToken, err := h.AuthUsecase.VerifyOTP(payload)
	if err != nil {
		return res.Error(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h AuthHandler) Login(ctx *fiber.Ctx) error {
	payload := new(dto.LoginRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	accessToken, refreshToken, err := h.AuthUsecase.Login(payload)
	if err != nil {
		return res.Error(ctx, err)
	}

	return ctx.JSON(fiber.Map{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
