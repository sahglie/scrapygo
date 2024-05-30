package services

import (
	"log/slog"
	"scrapygo/internal/config"
	"scrapygo/internal/database"
)

type Config struct {
	DB     *database.Queries
	Logger *slog.Logger
}

func NewServiceConfig() *Config {
	cfg := config.NewConfig()
	return &Config{
		DB:     cfg.DB,
		Logger: cfg.Logger,
	}
}

func NewServiceTestConfig() *Config {
	cfg := config.NewConfigTest()
	return &Config{
		DB:     cfg.DB,
		Logger: cfg.Logger,
	}
}
