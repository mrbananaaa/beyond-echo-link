package main

import (
	"fmt"
	"log"

	"github.com/mrbananaaa/bel-server/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	fmt.Println("AppEnv:", cfg.AppEnv)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("RedisAddr:", cfg.RedisAddr)
	fmt.Println("Debug:", cfg.Debug)

	fmt.Println("Wassup bbg")
}
