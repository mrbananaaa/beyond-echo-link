package apperror

import (
	"net/http"
)

func Internal(
	errType ErrorType,
	err error,
) *AppError {
	return New(
		errType,
		CodeInternal,
		"internal server error",
		http.StatusInternalServerError,
		err,
	)
}

func BadRequest(
	errType ErrorType,
	msg string,
	err error,
) *AppError {
	return New(
		errType,
		CodeBadRequest,
		msg,
		http.StatusBadRequest,
		err,
	)
}

func Validation(
	err error,
	field ...string,
) *AppError {
	return New(
		TypeInfrastructure,
		CodeBadRequest,
		"validation error",
		http.StatusBadRequest,
		err,
		field...,
	)
}

func Conflict(
	errType ErrorType,
	msg string,
	field string,
	err error,
) *AppError {
	return New(
		errType,
		CodeConflict,
		msg,
		http.StatusConflict,
		err,
		field,
	)
}

func Invalid(
	errType ErrorType,
	msg string,
	field string,
	err error,
) *AppError {
	return New(
		errType,
		CodeBadRequest,
		msg,
		http.StatusBadRequest,
		err,
		field,
	)
}

func InvalidCredentials(
	errType ErrorType,
	msg string,
	err error,
) *AppError {
	return New(
		errType,
		CodeUnauthorized,
		msg,
		http.StatusUnauthorized,
		err,
	)
}
