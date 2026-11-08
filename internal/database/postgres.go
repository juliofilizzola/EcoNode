package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"EcoNode/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(buildConnectionURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("parse database connection settings: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.ConnConfig.ConnectTimeout = cfg.ConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func buildConnectionURL(cfg config.DatabaseConfig) string {
	connectionURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Path:     "/" + cfg.Name,
		RawQuery: url.Values{"sslmode": {cfg.SSLMode}}.Encode(),
	}

	return connectionURL.String()
}
