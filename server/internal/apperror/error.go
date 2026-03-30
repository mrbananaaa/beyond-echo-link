package apperror

import (
	"errors"
	"net/http"
)

const (
	CodeInternal     = "INTERNAL_ERROR"
	CodeBadRequest   = "BAD_REQUEST"
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeConflict     = "CONFLICT"
)

var (
	ErrInternal     = New(TypeInfrastructure, CodeInternal, "internal server error", http.StatusInternalServerError, nil)
	ErrBadRequest   = New(TypeInfrastructure, CodeBadRequest, "bad request", http.StatusBadRequest, nil)
	ErrNotFound     = New(TypeInfrastructure, CodeNotFound, "resource not found", http.StatusNotFound, nil)
	ErrUnauthorized = New(TypeInfrastructure, CodeUnauthorized, "unauthorized", http.StatusUnauthorized, nil)
)

// AppEror internal error structure
type AppEror struct {
	Type    Type     // typeof error for logging purpose
	Code    string   // code message for client
	Status  int      // http status code
	Message string   // error message for client
	Details []string // fields or other detail messages for client
	Err     error    // original error
}

func (e *AppEror) Error() string {
	return e.Message
}

func (e *AppEror) Unwrap() error {
	return e.Err
}

func New(t Type, code, message string, status int, err error, details ...string) *AppEror {
	return &AppEror{
		Type:    t,
		Code:    code,
		Status:  status,
		Message: message,
		Details: details,
		Err:     err,
	}
}

func Typeof(err error) Type {
	if appErr, ok := errors.AsType[*AppEror](err); ok {
		return appErr.Type
	}

	return TypeInfrastructure
}

func ValidationError(field ...string) *AppEror {
	return New(
		TypeInfrastructure,
		CodeBadRequest,
		"validation error",
		http.StatusBadRequest,
		nil,
		field...,
	)
}

func Conflict(msg, field string, err error) *AppEror {
	return New(
		TypeDB,
		CodeConflict,
		msg,
		http.StatusConflict,
		err,
		field,
	)
}

func Invalid(msg, field string, err error) *AppEror {
	return New(
		TypeDB,
		CodeBadRequest,
		msg,
		http.StatusBadRequest,
		err,
		field,
	)
}

func Internal(err error) *AppEror {
	return New(
		TypeDB,
		CodeInternal,
		"internal server error",
		http.StatusInternalServerError,
		err,
	)
}

func InvalidCredentials(msg string, err error) *AppEror {
	return New(
		TypeBusiness,
		CodeUnauthorized,
		msg,
		http.StatusUnauthorized,
		err,
	)
}
