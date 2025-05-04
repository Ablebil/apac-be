package env

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv  string `env:"APP_ENV"`
	AppHost string `env:"APP_HOST"`
	AppPort string `env:"APP_PORT"`
	AppUrl  string `env:"APP_URL"`

	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`

	TokenLength int `env:"TOKEN_LENGTH"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectUrl  string `enc:"GOOGLE_REDIRECT_URL"`
}

func New() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	_env := new(Env)
	if err := env.Parse(_env); err != nil {
		return nil, err
	}

	return _env, nil
}
