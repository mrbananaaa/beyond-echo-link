package main

import (
	"fmt"

	"github.com/mrbananaaa/bel-server/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("AppEnv:", cfg.AppEnv)
	fmt.Println("Port:", cfg.Port)
	fmt.Println("RedisAddr:", cfg.RedisAddr)
	fmt.Println("Debug:", cfg.Debug)

	fmt.Println("Wassup bbg")
}
