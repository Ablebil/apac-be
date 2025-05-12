package rest

import (
	"apac/internal/app/user/usecase"
	"apac/internal/domain/dto"
	"apac/internal/infra/helper"
	res "apac/internal/infra/response"
	"apac/internal/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	Validator   *validator.Validate
	UserUsecase usecase.UserUsecaseItf
	helper      helper.HelperItf
}

func NewUserHandler(routerGroup fiber.Router, userUsecase usecase.UserUsecaseItf, validator *validator.Validate, m middleware.MiddlewareItf, helper helper.HelperItf) {
	UserHandler := UserHandler{
		Validator:   validator,
		UserUsecase: userUsecase,
		helper:      helper,
	}

	routerGroup = routerGroup.Group("/user")
	routerGroup.Get("/profile", m.Authentication, UserHandler.GetProfile)
	routerGroup.Patch("/profile", m.Authentication, UserHandler.EditProfile)
	routerGroup.Post("/preferences", m.Authentication, UserHandler.AddPreference)
	routerGroup.Delete("/preferences", m.Authentication, UserHandler.RemovePreference)
}

func (h UserHandler) GetProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userID").(uuid.UUID)

	user, err := h.UserUsecase.GetProfile(userId)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "User retrieved successfully", fiber.Map{
		"user": user,
	})
}

func (h UserHandler) EditProfile(ctx *fiber.Ctx) error {
	payload := new(dto.EditProfileRequest)
	if err := h.helper.FormParser(ctx, payload); err != nil {
		return res.BadRequest(ctx, err.Error())
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	userId := ctx.Locals("userID").(uuid.UUID)
	if err := h.UserUsecase.EditProfile(userId, payload); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Profile updated", nil)
}

func (h UserHandler) AddPreference(ctx *fiber.Ctx) error {
	payload := new(dto.AddPreferenceRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	userId := ctx.Locals("userID").(uuid.UUID)

	if err := h.UserUsecase.AddPreference(userId, payload); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Preference updated", nil)
}

func (h UserHandler) RemovePreference(ctx *fiber.Ctx) error {
	payload := new(dto.RemovePreferenceRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	userId := ctx.Locals("userID").(uuid.UUID)

	if err := h.UserUsecase.RemovePreference(userId, payload); err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "Preference removed", nil)
}
