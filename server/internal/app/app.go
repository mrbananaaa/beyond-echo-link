package app

import (
	"net/http"

	apphttp "github.com/mrbananaaa/bel-server/internal/http"
)

type App struct {
	server *apphttp.Server
}

func New() (*App, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	server := apphttp.New(":8080", mux)

	return &App{
		server: server,
	}, nil
}

func (a *App) Run() error {
	return a.server.Run()
}
