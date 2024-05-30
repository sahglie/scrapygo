package services

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"scrapygo/internal/database"
)

type Config struct {
	DB     *database.Queries
	Logger *slog.Logger
}

func NewServiceConfig(logger *slog.Logger) *Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	return &Config{
		DB:     database.New(db),
		Logger: logger,
	}
}

func NewServiceTestConfig(logger *slog.Logger) *Config {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		panic(err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}

	return &Config{
		DB:     database.New(db),
		Logger: logger,
	}
}
