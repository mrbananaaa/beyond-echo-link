package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	queries "github.com/mrbananaaa/bel-server/internal/db/sqlc"
)

func WithTx(
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(q *queries.Queries) error,
) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := queries.New(tx)

	if err := fn(q); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
