package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func New() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("passwd", func(f1 validator.FieldLevel) bool {
		password := f1.Field().String()

		hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
		hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
		hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

		return hasLower && hasUpper && hasSpecial && hasNumber
	})

	return v
}
