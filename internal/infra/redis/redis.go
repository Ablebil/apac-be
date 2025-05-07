package redis

import (
	"apac/internal/domain/env"
	"time"

	"github.com/gofiber/storage/redis"
)

type RedisItf interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte, exp time.Duration) error
	Delete(key string) error
	Reset() error
	Close() error
	SetOTP(email, otp string, exp time.Duration) error
	GetOTP(email string) (string, error)
	DeleteOTP(email string) error
}

type Redis struct {
	store *redis.Storage
}

func NewRedis(env *env.Env) RedisItf {
	return &Redis{
		store: redis.New(redis.Config{
			Host:     env.RedisHost,
			Port:     env.RedisPort,
			Username: env.RedisUsername,
			Password: env.RedisPassword,
		}),
	}
}

func (r *Redis) Get(key string) ([]byte, error) {
	val, err := r.store.Get(key)
	if err != nil {
		return make([]byte, 0), err
	}

	return val, err
}

func (r *Redis) Set(key string, val []byte, exp time.Duration) error {
	return r.store.Set(key, val, exp)
}

func (r *Redis) Delete(key string) error {
	return r.store.Delete(key)
}

func (r *Redis) Reset() error {
	return r.store.Reset()
}

func (r *Redis) Close() error {
	return r.store.Close()
}

func (r *Redis) SetOTP(email, otp string, exp time.Duration) error {
	key := "otp:" + email
	return r.Set(key, []byte(otp), exp)
}

func (r *Redis) GetOTP(email string) (string, error) {
	key := "otp:" + email
	val, err := r.Get(key)
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (r *Redis) DeleteOTP(email string) error {
	key := "otp:" + email
	return r.Delete(key)
}
