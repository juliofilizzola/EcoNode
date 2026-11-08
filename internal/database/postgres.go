package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"EcoNode/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnection(ctx context.Context, cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(buildConnectionURL(cfg)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("access sql database handle: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConns)

	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB, models ...any) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto migrate schema: %w", err)
	}

	return nil
}

func buildConnectionURL(cfg config.DatabaseConfig) string {
	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:   "/" + cfg.Name,
		RawQuery: url.Values{
			"sslmode":         {cfg.SSLMode},
			"connect_timeout": {strconv.Itoa(int(cfg.ConnectTimeout.Seconds()))},
		}.Encode(),
	}

	return connectionURL.String()
}
