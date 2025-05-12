package response

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (e *Err) Error() string {
	return e.Message
}

func newError(code int, defaultMsg string, message ...string) *Err {
	msg := defaultMsg
	if len(message) == 1 {

		msg = message[0]
	}
	return &Err{Code: code, Message: msg}
}

func ErrNotFound(message ...string) *Err {
	return newError(fiber.ErrNotFound.Code, fiber.ErrNotFound.Message, message...)
}

func ErrBadRequest(message ...string) *Err {
	return newError(fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, message...)

}

func ErrValidationError(err any) *Err {
	payload := map[string]any{"errors": err}
	return &Err{Code: fiber.ErrBadRequest.Code, Message: fiber.ErrBadRequest.Message, Payload: payload}
}

func ErrInternalServer(message ...string) *Err {
	return newError(fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message, message...)
}

func ErrUnauthorized(message ...string) *Err {
	return newError(fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, message...)
}

func ErrForbidden(message ...string) *Err {
	return newError(fiber.ErrForbidden.Code, fiber.ErrForbidden.Message, message...)
}

func ErrConflict(message ...string) *Err {
	return newError(fiber.ErrConflict.Code, fiber.ErrConflict.Message, message...)
}

func ErrUnprocessableEntity(message ...string) *Err {
	return newError(fiber.ErrUnprocessableEntity.Code, fiber.ErrUnprocessableEntity.Message, message...)
}

func ErrEntityTooLarge(message ...string) *Err {
	return newError(fiber.ErrRequestEntityTooLarge.Code, fiber.ErrRequestEntityTooLarge.Message, message...)
}

func respondWithError(ctx *fiber.Ctx, code int, defaultMsg string, message ...string) error {
	msg := defaultMsg
	if len(message) == 1 {
		msg = message[0]
	}
	return ctx.Status(code).JSON(Res{
		StatusCode: code,
		Message:    msg,
	})
}

func BadRequest(ctx *fiber.Ctx, message ...string) error {
	return respondWithError(ctx, fiber.ErrBadRequest.Code, fiber.ErrBadRequest.Message, message...)
}

func ValidationError(ctx *fiber.Ctx, err error) error {
	errorsMap := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		errorsMap[field] = strings.Trim(fmt.Sprintf("%s: %s %s", field, err.Tag(), err.Param()), " ")
	}

	payload := map[string]interface{}{"errors": errorsMap}

	return ctx.Status(fiber.ErrBadRequest.Code).JSON(Res{
		StatusCode: fiber.ErrBadRequest.Code,
		Message:    fiber.ErrBadRequest.Message,
		Payload:    payload,
	})
}

func InternalServerError(ctx *fiber.Ctx) error {
	return respondWithError(ctx, fiber.ErrInternalServerError.Code, fiber.ErrInternalServerError.Message)
}

func Unauthorized(ctx *fiber.Ctx, message ...string) error {
	return respondWithError(ctx, fiber.ErrUnauthorized.Code, fiber.ErrUnauthorized.Message, message...)
}

func Forbidden(ctx *fiber.Ctx, message ...string) error {
	return respondWithError(ctx, fiber.ErrForbidden.Code, fiber.ErrForbidden.Message, message...)
}

func Error(ctx *fiber.Ctx, err *Err) error {
	var customErr *Err
	if errors.As(err, &customErr) {
		return ctx.Status(customErr.Code).JSON(Res{
			StatusCode: customErr.Code,
			Message:    customErr.Message,
			Payload:    customErr.Payload,
		})
	}

	return InternalServerError(ctx)
}
