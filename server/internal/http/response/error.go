package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type ErrorResponse struct {
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	l := logger.FromContext(r.Context())

	var appErr *apperror.Error
	resp := &ErrorResponse{}

	if errors.As(err, &appErr) {
		if appErr.Code == apperror.CodeInternal {
			l.Error("internal error",
				"error", err.Error(),
				"code", appErr.Code,
			)
		} else {
			l.Info("client error",
				"error", err.Error(),
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

	l.Error("unexpected error",
		"error", err.Error(),
		"type", fmt.Sprintf("%T", err),
	)

	resp.Error = "INTERNAL_ERROR"
	resp.Message = "internal server error"

	JSON(w, http.StatusInternalServerError, resp)

}

func statusFromCode(code string) int {
	switch code {
	case apperror.CodeBadRequest:
		return http.StatusBadRequest
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
