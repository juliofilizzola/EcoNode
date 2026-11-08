package main

import (
	"context"
	"fmt"
	"log"

	"EcoNode/internal/config"
	"EcoNode/internal/database"
	"EcoNode/internal/route"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.ConnectTimeout)
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()

	engine := route.SetupEngine()

	return engine.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
