package services

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestConfig_ScrapeFeed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	config := NewServiceConfig(logger)

	feed, err := config.DB.GetFeedByUrl(context.TODO(), "https://blog.boot.dev/index.xml")
	assert.NoError(t, err)

	config.ScrapeFeed(feed)
}
