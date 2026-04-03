package validation

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mrbananaaa/bel-server/internal/apperror"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

func (v *Validator) Validate(i any) error {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	var errors []string

	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, formatError(err))
	}

	return apperror.New(
		apperror.TypeInfrastructure,
		apperror.CodeBadRequest,
		"validation error",
		http.StatusBadRequest,
		err,
		errors...,
	)
}

func formatError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "gte":
		return fmt.Sprintf("%s must be >= %s", err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s must be <= %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
