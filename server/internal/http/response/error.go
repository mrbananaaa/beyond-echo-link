package response

import (
	"errors"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/apperror"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type ErrResponse struct {
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	l := logger.FromContext(r.Context())

	var appErr *apperror.AppEror
	resp := ErrResponse{}

	if errors.As(err, &appErr) {
		if appErr.Code == apperror.CodeInternal {
			l.Error("internal error",
				"err", err,
				"code", appErr.Code,
			)
		} else {
			l.Info("client error",
				"err", err,
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

	// fallback
	l.Error("unexpected error",
		"err", err,
	)

	resp.Error = apperror.CodeInternal
	resp.Message = "internal server error"

	JSON(w, http.StatusInternalServerError, resp)
}

func statusFromCode(code string) int {
	switch code {
	case apperror.CodeInternal:
		return http.StatusInternalServerError
	case apperror.CodeBadRequest:
		return http.StatusBadRequest
	case apperror.CodeNotFound:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
