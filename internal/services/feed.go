package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"scrapygo/internal/database"
	"scrapygo/internal/scraper"
	"time"
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

var ErrPartialScrape = errors.New("partial scrape")
var ErrFailedScrape = errors.New("failed scrape")

func (cfg *Config) ScrapeFeed(feed database.Feed) error {
	cfg.Logger.Info("attempting to fetch feed", "url", feed.Url)

	feedData, err := scraper.FetchFeed(feed.Url)
	if err != nil {
		cfg.Logger.Error("failed to fetch feed", "err", err)
		return err
	}

	posts, err := cfg.DB.GetPostsByFeedID(context.TODO(), feed.ID)
	if err != nil {
		cfg.Logger.Error("failed to fetch feed", "err", err)
		return err
	}

	failedPostUrls := make([]string, 0)

	for _, p := range feedData.Posts {
		if postAlreadyScraped(p.Link, posts) {
			continue
		}

		_, err := cfg.DB.CreatePost(context.TODO(), database.CreatePostParams{
			ID:     uuid.New(),
			FeedID: feed.ID,
			Title:  p.Title,
			Description: sql.NullString{
				String: p.Description,
				Valid:  len(p.Description) > 0,
			},
			Url:         p.Link,
			PublishedAt: p.PubDate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})

		if err != nil {
			cfg.Logger.Error("failed to create post", "url", p.Link, "feedId", feed.ID, "err", err)
			failedPostUrls = append(failedPostUrls, p.Link)
		} else {
			cfg.Logger.Error("created post", "url", p.Link)
		}

	}

	if len(failedPostUrls) == 0 {
		return nil
	}

	if len(failedPostUrls) == len(feedData.Posts) {
		return ErrFailedScrape
	}

	return ErrPartialScrape
}

func postAlreadyScraped(postURL string, posts []database.Post) bool {
	for _, p := range posts {
		if p.Url == postURL {
			return true
		}
	}
	return false
}
