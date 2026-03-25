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
	return q.CreateUser(ctx, params)
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (queries.User, error) {
	q := r.getQueries(ctx)
	return q.GetUserByID(ctx, userID)
}

func (r *UserRepository) GetUserByLookupID(ctx context.Context, lookupID string) (queries.User, error) {
	q := r.getQueries(ctx)
	return q.GetUserByLookupID(ctx, lookupID)
}
