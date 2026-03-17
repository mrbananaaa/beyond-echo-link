package app

import (
	"context"
	"fmt"

	"github.com/mrbananaaa/bel-server/internal/config"
	apphttp "github.com/mrbananaaa/bel-server/internal/http"
	"github.com/mrbananaaa/bel-server/internal/http/handlers"
)

type App struct {
	Config *config.Config
	server *apphttp.Server
}

func New() (*App, error) {
	cfg := config.MustLoad()

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
		server: server,
	}, nil
}

func (a *App) Start() error {
	return a.server.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	// INFO: cleaning the app here...

	return a.server.Shutdown(ctx)
}
