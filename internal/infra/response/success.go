package response

import "github.com/gofiber/fiber/v2"

func SuccessResponse(ctx *fiber.Ctx, message string, data any) error {
	return ctx.Status(fiber.StatusOK).JSON(Res{
		StatusCode: fiber.StatusOK,
		Message:    message,
		Payload:    data,
	})
}
