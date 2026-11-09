package main

import (
	"context"
	"fmt"
	"log"

	"EcoNode/internal/config"
	"EcoNode/internal/controller"
	"EcoNode/internal/database"
	"EcoNode/internal/middleware"
	"EcoNode/internal/repository"
	"EcoNode/internal/route"
	"EcoNode/internal/service"
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

	db, err := database.NewPostgresConnection(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	if err := database.AutoMigrate(db, &repository.User{}); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("access sql database handle: %w", err)
	}
	defer sqlDB.Close()

	userRepository := repository.NewUserRepository(db)
	tokenService := service.NewJWTTokenService(cfg.Auth.JWTSecret, cfg.Auth.JWTTTL)
	authService := service.NewAuthService(userRepository, tokenService)
	authController := controller.NewAuthController(authService)
	authenticate := middleware.Authenticate(tokenService)

	engine := route.SetupEngine(authController, authenticate)

	return engine.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
