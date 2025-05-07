package rest

import (
	"apac/internal/app/gemini/usecase"
	"apac/internal/domain/dto"
	res "apac/internal/infra/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GeminiHandler struct {
	Validator     *validator.Validate
	GeminiUsecase usecase.GeminiUsecaseItf
}

func NewGeminiHandler(routerGroup fiber.Router, geminiUsecase usecase.GeminiUsecaseItf, validator *validator.Validate) {
	geminiHandler := GeminiHandler{
		Validator:     validator,
		GeminiUsecase: geminiUsecase,
	}

	routerGroup = routerGroup.Group("/gemini")
	routerGroup.Post("/", geminiHandler.Prompt)
}

func (h GeminiHandler) Prompt(ctx *fiber.Ctx) error {
	payload := new(dto.GeminiRequest)
	if err := ctx.BodyParser(&payload); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(payload); err != nil {
		return res.ValidationError(ctx, err)
	}

	response, err := h.GeminiUsecase.Prompt(payload)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, "AI prompt succesful", response)
}
