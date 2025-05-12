package middleware

import (
	"apac/internal/infra/jwt"

	"github.com/gofiber/fiber/v2"
)

type MiddlewareItf interface {
	Authentication(*fiber.Ctx) error
}

type Middleware struct {
	jwt jwt.JWTItf
}

func New(jwt jwt.JWTItf) MiddlewareItf {
	return &Middleware{
		jwt: jwt,
	}
}
