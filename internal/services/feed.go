package services

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"scrapygo/internal/database"
	"scrapygo/internal/scraper"
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

// Take a url of a feed, and load any new entries into
// feeds table
func (cfg *Config) LoadFeed(url string) error {
	cfg.Logger.Info("attempting to fetch feed", "url", url)

	feed, err := scraper.FetchFeed(url)
	if err != nil {
		cfg.Logger.Error("failed to fetch feed", "err", err)
	}

	feed
	// for each feed entry, check to see if it is in the db. If not in db
	// add it.

	// if in db and timestamps are equal, ignore it.

	// if in db and timestamps are not equal, update it.

	return nil
}
