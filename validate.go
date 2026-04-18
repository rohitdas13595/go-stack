package gostack

import "github.com/go-playground/validator/v10"

type validateWrapper struct {
	v *validator.Validate
}

func newValidateWrapper() *validateWrapper {
	return &validateWrapper{v: validator.New()}
}
