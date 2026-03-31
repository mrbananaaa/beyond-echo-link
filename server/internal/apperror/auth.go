package apperror

import "net/http"

func InvalidCredentials(t Type, msg string, err error) *AppEror {
	return New(
		TypeBusiness,
		CodeUnauthorized,
		msg,
		http.StatusUnauthorized,
		err,
	)
}
