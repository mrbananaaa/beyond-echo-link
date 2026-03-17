package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("server started on :%v", a.Config.Server.Port)

		if err := a.server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Printf("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	return a.server.Shutdown(shutdownCtx)
}
