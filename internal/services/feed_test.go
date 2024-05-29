package services

import (
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"testing"
)

func TestConfig_LoadFeed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	url := "https://blog.boot.dev/index.xml"

	config := NewServiceConfig(logger)
	config.LoadFeed(url)
}
