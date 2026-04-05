package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/bel-server/internal/infra/db"
	queries "github.com/mrbananaaa/bel-server/internal/infra/db/sqlc"
)

type UserRepository interface {
	CreateUser(ctx context.Context, params queries.CreateUserParams) (queries.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (queries.User, error)
	GetUserByLookupID(ctx context.Context, lookupID string) (queries.User, error)
	GetUserByUsername(ctx context.Context, username string) (queries.User, error)
}

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{
		pool: pool,
	}
}

func (r *postgresUserRepository) getQueries(ctx context.Context) *queries.Queries {
	if tx, ok := db.ExtractTx(ctx); ok {
		return queries.New(tx)
	}
	return queries.New(r.pool)
}

func (r *postgresUserRepository) CreateUser(ctx context.Context, params queries.CreateUserParams) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.CreateUser(ctx, params)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}

func (r *postgresUserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByID(ctx, userID)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}

func (r *postgresUserRepository) GetUserByLookupID(ctx context.Context, lookupID string) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByLookupID(ctx, lookupID)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}

func (r *postgresUserRepository) GetUserByUsername(ctx context.Context, username string) (queries.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByUsername(ctx, username)
	if err != nil {
		return queries.User{}, db.MapError(err)
	}

	return u, nil
}
