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
	authHandler := AuthHandler{
		Validator:   validator,
		AuthUsecase: authUsecase,
	}

	routerGroup = routerGroup.Group("/auth")
	routerGroup.Post("/register", authHandler.Register)
	routerGroup.Post("/verify-otp", authHandler.VerifyOTP)
	routerGroup.Post("/choose-preference", authHandler.ChoosePreference)
	routerGroup.Post("/login", authHandler.Login)
	routerGroup.Post("/refresh-token", authHandler.RefreshToken)
	routerGroup.Post("/logout", authHandler.Logout)
	routerGroup.Get("/google", authHandler.GoogleLogin)
	routerGroup.Get("/google/callback", authHandler.GoogleCallback)
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

	return res.SuccessResponse(ctx, "Verification successful", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
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

	return res.SuccessResponse(ctx, "Login successful", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	payload := new(dto.RefreshToken)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	accessToken, refreshToken, err := h.AuthUsecase.RefreshToken(payload)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Token refreshed", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h AuthHandler) Logout(ctx *fiber.Ctx) error {
	payload := new(dto.LogoutRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.Logout(payload); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Logout successful", nil)
}

func (h AuthHandler) GoogleLogin(ctx *fiber.Ctx) error {
	url, err := h.AuthUsecase.GoogleLogin()
	if err != nil {
		return res.Error(ctx, err)
	}

	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h AuthHandler) GoogleCallback(ctx *fiber.Ctx) error {
	payload := &dto.GoogleCallbackRequest{
		Code:  ctx.Query("code"),
		State: ctx.Query("state"),
		Error: ctx.Query("error"),
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	accessToken, refreshToken, err := h.AuthUsecase.GoogleCallback(payload)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Token refreshed", fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h AuthHandler) ChoosePreference(ctx *fiber.Ctx) error {
	payload := new(dto.ChoosePreferenceRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	if err := h.AuthUsecase.ChoosePreference(payload); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Preference updated", nil)
}
