package env

import (
	"time"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Env struct {
	AppEnv  string `env:"APP_ENV"`
	AppHost string `env:"APP_HOST"`
	AppPort int    `env:"APP_PORT"`
	AppUrl  string `env:"APP_URL"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUsername string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisUsername string `env:"REDIS_USER"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectUrl  string `env:"GOOGLE_REDIRECT_URL"`
	RedirectUrl        string `enc:"REDIRECT_URL"`

	EmailUser string `env:"EMAIL_USER"`
	EmailPass string `env:"EMAIL_PASS"`

	AccessSecret  string `env:"JWT_SECRET"`
	RefreshSecret string `env:"JWT_REFRESH_SECRET"`

	StateLength int           `env:"STATE_LENGTH"`
	StateExpiry time.Duration `env:"STATE_EXPIRY"`

	SupabaseBucket  string `env:"SUPABASE_BUCKET"`
	SupabaseUrl     string `env:"SUPABASE_URL"`
	SupabaseAnonKey string `env:"SUPABASE_ANON_KEY"`

	DefaultProfilePic string `env:"DEFAULT_PROFILE_PIC"`

	GeminiAPIKey string `env:"GEMINI_API_KEY"`
	GeminiModel  string `env:"GEMINI_MODEL"`
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
