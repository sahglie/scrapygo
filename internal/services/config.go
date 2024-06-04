package services

import (
	"log/slog"
	"scrapygo/internal/database"
)

type Config struct {
	DB     *database.Queries
	Logger *slog.Logger
}

func NewConfig(db *database.Queries, logger *slog.Logger) *Config {
	return &Config{
		DB:     db,
		Logger: logger,
	}
}
