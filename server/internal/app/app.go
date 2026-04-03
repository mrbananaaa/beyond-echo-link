package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrbananaaa/bel-server/internal/infra/config"
	"github.com/mrbananaaa/bel-server/internal/infra/db"
	apphttp "github.com/mrbananaaa/bel-server/internal/infra/http"
	"github.com/mrbananaaa/bel-server/internal/infra/http/handlers"
	authHandler "github.com/mrbananaaa/bel-server/internal/infra/http/handlers/auth"
	"github.com/mrbananaaa/bel-server/internal/infra/http/middlewares"
	"github.com/mrbananaaa/bel-server/internal/logger"
	"github.com/mrbananaaa/bel-server/internal/usecase/auth"
	"github.com/mrbananaaa/bel-server/internal/usecase/user"
	"github.com/mrbananaaa/bel-server/internal/validation"
)

type App struct {
	Config *config.Config
	Log    *slog.Logger
	server *apphttp.Server
	dbpool *pgxpool.Pool
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

	txManager := db.NewTxManager(dbpool)

	userRepo := user.NewUserRepository(dbpool)

	authService := auth.NewAuthService(
		txManager,
		userRepo,
		log,
		cfg.Server.JwtSecret,
		"bel-backend-dev",
		15*time.Minute,
	)

	validator := validation.New()
	healthHandler := handlers.NewHealthHandler()
	authHandler := authHandler.NewAuthHandler(validator, authService)

	logMiddleware := middlewares.NewLogMiddleware(log)
	authMiddleware := middlewares.NewAuthMiddleware(authService, log)

	router := apphttp.NewRouter(
		apphttp.Handlers{
			Health: healthHandler,
			Auth:   authHandler,
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

	return a.server.Shutdown(ctx)
}
