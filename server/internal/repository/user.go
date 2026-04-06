package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/bel-server/internal/domain/user"
	"github.com/mrbananaaa/bel-server/internal/infra/db"
	queries "github.com/mrbananaaa/bel-server/internal/infra/db/sqlc"
)

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) user.UserRepository {
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

func (r *postgresUserRepository) Save(ctx context.Context, user *user.User) (*user.User, error) {
	q := r.getQueries(ctx)

	u, err := q.CreateUser(ctx, queries.CreateUserParams{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		LookupID: user.LookupID,
		Bio: pgtype.Text{
			String: user.Bio,
		},
		ProfilePicture: pgtype.Text{
			String: user.ProfilePicture,
		},
	})
	if err != nil {
		return nil, db.MapError(err)
	}

	user.FromDB(u)

	return user, nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByID(ctx, userID)
	if err != nil {
		return nil, db.MapError(err)
	}
	user := &user.User{}
	user.FromDB(u)

	return user, nil
}

func (r *postgresUserRepository) GetByLookupID(ctx context.Context, lookupID string) (*user.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByLookupID(ctx, lookupID)
	if err != nil {
		return nil, db.MapError(err)
	}
	user := &user.User{}
	user.FromDB(u)

	return user, nil
}

func (r *postgresUserRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	q := r.getQueries(ctx)

	u, err := q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, db.MapError(err)
	}
	user := &user.User{}
	user.FromDB(u)

	return user, nil
}
