package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrbananaaa/bel-server/internal/domain/apperror"
)

func MapError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return mapUniqueViolation(pgErr)
		case "23503": // foreign_key_violation
			return apperror.Invalid(
				apperror.TypeDB,
				"invalid reference",
				"",
				err,
			)
		case "23502": // not_null_violation
			return apperror.Invalid(
				apperror.TypeDB,
				"missing required field",
				pgErr.ColumnName,
				err,
			)

		default:
			return apperror.Internal(
				apperror.TypeDB,
				err,
			)
		}
	}

	return err
}

func mapUniqueViolation(pgErr *pgconn.PgError) error {
	switch pgErr.ConstraintName {
	case "users_email_key":
		return apperror.Conflict(
			apperror.TypeDB,
			"email already exists",
			"email",
			pgErr,
		)
	case "users_username_key":
		return apperror.Conflict(
			apperror.TypeDB,
			"username already exists",
			"username",
			pgErr,
		)
	case "users_lookup_id_key":
		return apperror.Conflict(
			apperror.TypeDB,
			"lookup_id already exists",
			"username",
			pgErr,
		)
	default:
		return apperror.Conflict(
			apperror.TypeDB,
			"resource already exists",
			"",
			pgErr,
		)
	}
}
