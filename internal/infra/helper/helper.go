package helper

import (
	"apac/internal/domain/env"
	res "apac/internal/infra/response"
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type HelperItf interface {
	ValidateImage(file *multipart.FileHeader) *res.Err
	FormParser(ctx *fiber.Ctx, target interface{}) error
}

type Helper struct {
	env *env.Env
}

func NewHelper(env *env.Env) HelperItf {
	return &Helper{env: env}
}

func (h Helper) FormParser(ctx *fiber.Ctx, target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)
		formKey := structField.Tag.Get("form")

		if formKey == "" || !field.CanSet() {
			continue
		}

		if values, ok := form.Value[formKey]; ok && len(values) > 0 {
			valueStr := values[0]

			switch field.Kind() {
			case reflect.String:
				field.SetString(valueStr)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, _ := strconv.ParseInt(valueStr, 10, 64)
				field.SetInt(intValue)
			case reflect.Float32, reflect.Float64:
				floatValue, _ := strconv.ParseFloat(valueStr, 64)
				field.SetFloat(floatValue)
			case reflect.Bool:
				boolValue, _ := strconv.ParseBool(valueStr)
				field.SetBool(boolValue)
			}
		}

		if files, ok := form.File[formKey]; ok && len(files) > 0 {
			if field.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
				field.Set(reflect.ValueOf(files[0]))
			} else if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
				field.Set(reflect.ValueOf(files))
			}
		}
	}

	return nil
}

func (h Helper) ValidateImage(file *multipart.FileHeader) *res.Err {
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return res.ErrUnprocessableEntity("Invalid file type")
	}

	if file.Size > 10*1024*1024 {
		return res.ErrEntityTooLarge("Photo size must be less than 10MB")
	}

	return nil
}
