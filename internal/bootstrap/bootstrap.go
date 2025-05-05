package bootstrap

import (
	"apac/internal/domain/env"
	"apac/internal/infra/email"
	"apac/internal/infra/fiber"
	"apac/internal/infra/jwt"
	"apac/internal/infra/postgresql"
	"fmt"

	AuthHandler "apac/internal/app/auth/interface/rest"
	AuthRepo "apac/internal/app/auth/repository"
	AuthUsecase "apac/internal/app/auth/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Start() error {
	config, err := env.New()
	if err != nil {
		panic(err)
	}

	db, err := postgresql.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		config.DBHost,
		config.DBUsername,
		config.DBPassword,
		config.DBName,
		config.DBPort,
	), config)

	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}

	if err := postgresql.Migrate(db); err != nil {
		return err
	}

	v := validator.New()
	j := jwt.NewJWT(config)
	e := email.NewEmailService(config)

	app := fiber.New(config)
	app.Get("/metrics", monitor.New())
	v1 := app.Group("/api/v1")

	authRepository := AuthRepo.NewAuthRepository(db)

	authUsecase := AuthUsecase.NewAuthUsecase(config, db, authRepository, j, e)
	AuthHandler.NewAuthHandler(v1, authUsecase, v)

	return app.Listen(fmt.Sprintf("%s: %d", config.AppHost, config.AppPort))
}
