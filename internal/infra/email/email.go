package email

import (
	"apac/internal/domain/env"
	"fmt"

	"gopkg.in/gomail.v2"
)

type EmailItf interface {
	SendOTPEmail(to string, otp string) error
}

type Email struct {
	sender   string
	password string
}

func NewEmailService(env *env.Env) EmailItf {
	return &Email{
		sender:   env.EmailUser,
		password: env.EmailPass,
	}
}

func (e *Email) SendOTPEmail(to string, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.sender)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your OTP Code")
	m.SetBody("text/plain", fmt.Sprintf("Your OTP is: %s", otp))

	d := gomail.NewDialer("smtp.gmail.com", 587, e.sender, e.password)
	return d.DialAndSend(m)
}
