package bootstrap

import (
	"apac/internal/domain/env"
	"apac/internal/infra/fiber"
	"apac/internal/infra/postgresql"
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		panic(err)
	}

	db, err := postgresql.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.DBHost,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	), config)

	app := fiber.New(config)
	app.Get("/metrics", monitor.New())
	v1 := app.Group("/api/v1")

	return app.Listen(fmt.Sprintf("%s: %d", config.AppHost, config.AppPort))
}
