package service

import "github.com/go-playground/validator/v10"

type BoValidator interface {
	Validate(item any) error
}

type AccountValidator struct {
	validate *validator.Validate
}

func NewAccountValidator() BoValidator {
	return &AccountValidator{validate: validator.New()}
}

func (a *AccountValidator) Validate(item any) error {
	return a.validate.Struct(item)
}
