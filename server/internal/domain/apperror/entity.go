package apperror

import "errors"

type ErrorType string

const (
	TypeValidation     ErrorType = "validation_error"
	TypeBusiness       ErrorType = "business_error"
	TypeInfrastructure ErrorType = "infrastructure_error"
	TypeExternal       ErrorType = "externall_error"
	TypeDB             ErrorType = "db_error"
	TypeUnknown        ErrorType = "unknown_error"
)

func Typeof(err error) ErrorType {
	if appErr, ok := errors.AsType[*AppError](err); ok {
		return appErr.Type
	}
	return TypeUnknown
}

type ErrorCode string

const (
	CodeInternal     ErrorCode = "INTERNAL_ERROR"
	CodeBadRequest   ErrorCode = "BAD_REQUEST"
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeConflict     ErrorCode = "CONFLICT"
)

type AppError struct {
	Type    ErrorType
	Code    ErrorCode
	Status  int
	Message string
	Details []string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(
	errType ErrorType,
	errCode ErrorCode,
	message string,
	status int,
	err error,
	details ...string,
) *AppError {
	return &AppError{
		Type:    errType,
		Code:    errCode,
		Status:  status,
		Message: message,
		Details: details,
		Err:     err,
	}
}
