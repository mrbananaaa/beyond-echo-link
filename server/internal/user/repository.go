package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/bel-server/internal/db"
	queries "github.com/mrbananaaa/bel-server/internal/db/sqlc"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) getQueries(ctx context.Context) *queries.Queries {
	if tx, ok := db.ExtractTx(ctx); ok {
		return queries.New(tx)
	}
	return queries.New(r.pool)
}

func (r *UserRepository) CreateUser(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.CreateUser(ctx, params)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByID(ctx, userID)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}

func (r *UserRepository) GetUserByLookupID(ctx context.Context, lookupID string) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByLookupID(ctx, lookupID)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}
