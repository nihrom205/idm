package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()
	return &Validator{validate: validate}
}

func (v *Validator) Validate(request any) error {
	err := v.validate.Struct(request)
	if err != nil {
		var validateErr validator.ValidationErrors
		if errors.As(err, &validateErr) {
			return validateErr
		}
	}
	return err
}
