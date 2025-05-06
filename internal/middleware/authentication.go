package middleware

import (
	res "apac/internal/infra/response"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	auth := ctx.GetReqHeaders()["Authorization"]

	if len(auth) < 1 {
		return res.Unauthorized(ctx, "Missing access token")
	}

	token := strings.Split(auth[0], " ")

	if len(token) != 2 {
		return res.Unauthorized(ctx, "Invalid access token")
	}

	if token[0] != "Bearer" {
		return res.Unauthorized(ctx, "Wrong authorization type")
	}

	userID, name, email, err := m.jwt.VerifyAccessToken(token[1])
	if err != nil {
		return res.Unauthorized(ctx, err.Error())
	}

	ctx.Locals("userID", userID)
	ctx.Locals("name", name)
	ctx.Locals("email", email)

	return ctx.Next()
}
