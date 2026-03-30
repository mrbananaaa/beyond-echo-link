package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrbananaaa/bel-server/internal/apperror"
)

func MapError(err error) error {
	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		default:
			return apperror.ErrInternal
		}
	}

	return err
}

func mapUniqueViolation(pgErr *pgconn.PgError) error {
	switch pgErr.ConstraintName {

	}

	return nil
}
