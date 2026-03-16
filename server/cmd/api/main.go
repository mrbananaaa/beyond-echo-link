package main

import (
	"fmt"

	"github.com/mrbananaaa/bel-server/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println("AppEnv:", cfg.AppEnv)
	fmt.Println("Port:", cfg.Server.Port)
	fmt.Println("DbUrl:", cfg.DB.URL)
	fmt.Println("LogLevel:", cfg.Log.Level)

	fmt.Println("Wassup bbg")
}
