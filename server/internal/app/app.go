package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mrbananaaa/bel-server/internal/config"
	"github.com/mrbananaaa/bel-server/internal/db"
	apphttp "github.com/mrbananaaa/bel-server/internal/http"
	"github.com/mrbananaaa/bel-server/internal/http/handlers"
	"github.com/mrbananaaa/bel-server/internal/logger"
)

type App struct {
	Config *config.Config
	Log    *slog.Logger
	server *apphttp.Server
}

func New() (*App, error) {
	cfg := config.MustLoad()
	log := logger.New(logger.Config{
		Env:     "dev",
		Service: "infra",
	})

	dbpool, err := db.NewPool(cfg.DB.URL)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Error("cannot connect to db", "error", err.Error())
	}

	log.Info("connected to database")

	healthHandler := handlers.NewHealthHandler()

	router := apphttp.NewRouter(apphttp.Handlers{
		Health: healthHandler,
	})

	server := apphttp.NewServer(
		fmt.Sprintf(":%v", cfg.Server.Port),
		router,
	)

	return &App{
		Config: cfg,
		Log:    log,
		server: server,
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

	return a.server.Shutdown(ctx)
}
