package domain

import (
	"errors"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

var ErrUserNotFound = errors.New("user with such credentials not found")

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SignUpInput struct {
	Login    string `json:"login" validate:"required,gte=4"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=5"`
}
type SignInInput struct {
	Login    string `json:"login" validate:"required,gte=4"`
	Password string `json:"password" validate:"required,gte=5"`
}

func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}

func (i SignInInput) Validate() error {
	return validate.Struct(i)
}
