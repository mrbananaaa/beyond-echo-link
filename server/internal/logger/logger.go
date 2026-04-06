package logger

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/Marlliton/slogpretty"
	"github.com/mrbananaaa/bel-server/internal/domain/apperror"
)

type Config struct {
	Env     string
	Service string
}

func New(cfg Config) *slog.Logger {
	var handler slog.Handler

	if cfg.Env == "dev" {
		handler = slogpretty.New(os.Stdout, &slogpretty.Options{
			Level:      slog.LevelDebug,
			Colorful:   true,
			TimeFormat: slogpretty.DefaultTimeFormat,
			// AddSource: true,
			// Multiline:  true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	base := slog.New(handler)

	return base.With(
		"service", cfg.Service,
		"env", cfg.Env,
	)
}

// TODO: ADD LOG TYPE FOR FILTERING
func InfoEvent(l *slog.Logger, event string, msg string, args ...any) {
	l.Info(msg,
		append([]any{"event", event}, args...)...,
	)
}

func ErrorEvent(l *slog.Logger, event string, msg string, err error, args ...any) {
	if err == nil {
		l.Error(
			msg,
			"event", event,
			"error", errors.New("unknown error/nil").Error(),
			"error_type", "Unknown Error/nil",
		)
		return
	}

	errorType := apperror.Typeof(err)

	allArgs := []any{
		"event", event,
		"error", err.Error(),
		"error_type", string(errorType),
	}

	allArgs = append(allArgs, args...)

	l.Error(msg, allArgs...)
}

func DebugError(l *slog.Logger, msg string, err error) {
	if appErr, ok := errors.AsType[*apperror.AppError](err); ok {
		l.Debug(
			msg,
			"type", appErr.Type,
			"code", appErr.Code,
			"status", appErr.Status,
			"message", appErr.Message,
			"details", appErr.Details,
			"err", appErr.Err.Error(),
		)
		return
	}

	l.Debug(msg, "error", err)
}

func ErrorParseJSON(l *slog.Logger, err error) {
	ErrorEvent(l,
		"json_parsing_failed",
		"failed to decode request body",
		err,
	)
}

func ErrorValidation(l *slog.Logger, err error) {
	ErrorEvent(l,
		"req_validation_error",
		"failed to validate request body",
		apperror.New(
			apperror.TypeInfrastructure,
			apperror.CodeBadRequest,
			"validation error",
			http.StatusBadRequest,
			nil,
		),
	)
}
