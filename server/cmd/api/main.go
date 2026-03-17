package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrbananaaa/bel-server/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("server started on :%v", a.Config.Server.Port)
		if err := a.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}
