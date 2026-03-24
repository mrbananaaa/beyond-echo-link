package response

import (
	"errors"
<<<<<<< HEAD
	"fmt"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type ErrorResponse struct {
=======
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type ErrResponse struct {
>>>>>>> 35fc4f5 (feat: implementing apperror on error response)
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	l := logger.FromContext(r.Context())

<<<<<<< HEAD
	var appErr *apperror.Error
	resp := &ErrorResponse{}
=======
	var appErr *apperror.AppEror
	resp := ErrResponse{}
>>>>>>> 35fc4f5 (feat: implementing apperror on error response)

	if errors.As(err, &appErr) {
		if appErr.Code == apperror.CodeInternal {
			l.Error("internal error",
<<<<<<< HEAD
				"error", err.Error(),
=======
				"err", err,
>>>>>>> 35fc4f5 (feat: implementing apperror on error response)
				"code", appErr.Code,
			)
		} else {
			l.Info("client error",
<<<<<<< HEAD
				"error", err.Error(),
=======
				"err", err,
>>>>>>> 35fc4f5 (feat: implementing apperror on error response)
				"code", appErr.Code,
			)
		}

		resp.Error = appErr.Code
		resp.Message = appErr.Message

		if len(appErr.Details) > 0 {
			resp.Details = appErr.Details
		}

		JSON(w, statusFromCode(appErr.Code), resp)
		return
	}

<<<<<<< HEAD
	l.Error("unexpected error",
		"error", err.Error(),
		"type", fmt.Sprintf("%T", err),
	)

	resp.Error = "INTERNAL_ERROR"
	resp.Message = "internal server error"

	JSON(w, http.StatusInternalServerError, resp)

=======
	// fallback
	l.Error("unexpected error",
		"err", err,
	)

	resp.Error = apperror.CodeInternal
	resp.Message = "internal server error"

	JSON(w, http.StatusInternalServerError, resp)
>>>>>>> 35fc4f5 (feat: implementing apperror on error response)
}

func statusFromCode(code string) int {
	switch code {
<<<<<<< HEAD
	case apperror.CodeBadRequest:
		return http.StatusBadRequest
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
=======
	case apperror.CodeInternal:
		return http.StatusInternalServerError
	case apperror.CodeBadRequest:
		return http.StatusBadRequest
	case apperror.CodeNotFound:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
>>>>>>> 35fc4f5 (feat: implementing apperror on error response)
	default:
		return http.StatusInternalServerError
	}
}
