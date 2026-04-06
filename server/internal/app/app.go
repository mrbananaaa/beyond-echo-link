package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/bel-server/internal/domain/auth"
	"github.com/mrbananaaa/bel-server/internal/domain/token"
	"github.com/mrbananaaa/bel-server/internal/infra/config"
	"github.com/mrbananaaa/bel-server/internal/infra/db"
	apphttp "github.com/mrbananaaa/bel-server/internal/infra/http"
	"github.com/mrbananaaa/bel-server/internal/infra/http/handlers"
	authHandler "github.com/mrbananaaa/bel-server/internal/infra/http/handlers/auth"
	wsHandler "github.com/mrbananaaa/bel-server/internal/infra/http/handlers/ws"
	"github.com/mrbananaaa/bel-server/internal/infra/http/middlewares"
	redisinfra "github.com/mrbananaaa/bel-server/internal/infra/redis"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/repository"
	"github.com/mrbananaaa/bel-server/internal/validation"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config *config.Config
	Log    *slog.Logger
	server *apphttp.Server
	dbpool *pgxpool.Pool
	rdb    *redisinfra.Redis
}

func New() (*App, error) {
	cfg := config.MustLoad()
	log := logger.New(logger.Config{
		Env:     "dev",
		Service: "api",
	})

	dbpool, err := db.NewPool(cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Error("cannot connect to db", "error", err.Error())
		return nil, err
	}
	log.Info("connected to database")

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Address,
	})

	rdb := redisinfra.New(
		client,
		"dev",
		2*time.Second,
	)

	if err := rdb.Ping(context.Background()); err != nil {
		logger.ErrorEvent(log,
			"rdb_ping_failed",
			"failed to ping redis client",
			err,
		)
		return nil, err
	}

	tokenStore := redisinfra.NewTokenStore(rdb)

	txManager := db.NewTxManager(dbpool)

	userRepo := repository.NewUserRepository(dbpool)

	authService := auth.NewAuthService(
		txManager,
		userRepo,
		log,
	)
	tokenService := token.NewTokenService(
		tokenStore,
		cfg.Server.JwtSecret,
		"bel-backend-dev",
		15*time.Minute,
	)

	validator := validation.New()
	healthHandler := handlers.NewHealthHandler()
	authHandler := authHandler.NewAuthHandler(validator, authService, tokenService)
	wsHandler := wsHandler.NewWsHandler()

	logMiddleware := middlewares.NewLogMiddleware(log)
	authMiddleware := middlewares.NewAuthMiddleware(tokenService, log)

	router := apphttp.NewRouter(
		apphttp.Handlers{
			Health: healthHandler,
			Auth:   authHandler,
			Ws:     wsHandler,
		},
		apphttp.Middlewares{
			Log:  logMiddleware,
			Auth: authMiddleware,
		},
	)

	server := apphttp.NewServer(
		fmt.Sprintf(":%v", cfg.Server.Port),
		router,
	)

	return &App{
		Config: cfg,
		Log:    log,
		server: server,
		dbpool: dbpool,
		rdb:    rdb,
	}, nil
}

func (a *App) Start() error {
	if err := a.server.Start(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	// INFO: cleaning the app here...
	a.dbpool.Close()

	if err := a.rdb.Close(); err != nil {
		logger.ErrorEvent(a.Log,
			"rdb_closing_failed",
			"failed to close redis client",
			err,
		)
		return err
	}

	return a.server.Shutdown(ctx)
}
