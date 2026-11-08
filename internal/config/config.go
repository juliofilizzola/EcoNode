package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultServerPort             = 8080
	defaultDatabasePort           = 5432
	defaultDatabaseSSLMode        = "disable"
	defaultDatabaseMaxConns       = 10
	defaultDatabaseConnectTimeout = 5 * time.Second
	maxPortNumber                 = 65535
)

var validDatabaseSSLModes = []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Host           string
	Port           int
	User           string
	Password       string
	Name           string
	SSLMode        string
	MaxConns       int
	ConnectTimeout time.Duration
}

type envReader struct {
	errs []error
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("load environment file: %w", err)
	}

	reader := &envReader{}

	cfg := Config{
		Server: ServerConfig{
			Port: reader.port("SERVER_PORT", defaultServerPort),
		},
		Database: DatabaseConfig{
			Host:           reader.required("DB_HOST"),
			Port:           reader.port("DB_PORT", defaultDatabasePort),
			User:           reader.required("DB_USER"),
			Password:       reader.required("DB_PASSWORD"),
			Name:           reader.required("DB_NAME"),
			SSLMode:        reader.choice("DB_SSL_MODE", defaultDatabaseSSLMode, validDatabaseSSLModes),
			MaxConns:       reader.positiveInt("DB_MAX_CONNS", defaultDatabaseMaxConns),
			ConnectTimeout: reader.positiveDuration("DB_CONNECT_TIMEOUT", defaultDatabaseConnectTimeout),
		},
	}

	if err := errors.Join(reader.errs...); err != nil {
		return Config{}, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (r *envReader) required(key string) string {
	value := os.Getenv(key)
	if value == "" {
		r.errs = append(r.errs, fmt.Errorf("%s is required", key))
	}

	return value
}

func (r *envReader) withFallback(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func (r *envReader) choice(key, fallback string, options []string) string {
	value := r.withFallback(key, fallback)
	if !slices.Contains(options, value) {
		r.errs = append(r.errs, fmt.Errorf("%s must be one of %v, got %q", key, options, value))
	}

	return value
}

func (r *envReader) integer(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		r.errs = append(r.errs, fmt.Errorf("%s must be an integer: %w", key, err))
		return fallback
	}

	return value
}

func (r *envReader) port(key string, fallback int) int {
	value := r.integer(key, fallback)
	if value < 1 || value > maxPortNumber {
		r.errs = append(r.errs, fmt.Errorf("%s must be between 1 and %d, got %d", key, maxPortNumber, value))
	}

	return value
}

func (r *envReader) positiveInt(key string, fallback int) int {
	value := r.integer(key, fallback)
	if value < 1 {
		r.errs = append(r.errs, fmt.Errorf("%s must be greater than zero, got %d", key, value))
	}

	return value
}

func (r *envReader) positiveDuration(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		r.errs = append(r.errs, fmt.Errorf("%s must be a valid duration: %w", key, err))
		return fallback
	}

	if value <= 0 {
		r.errs = append(r.errs, fmt.Errorf("%s must be greater than zero, got %s", key, value))
	}

	return value
}
