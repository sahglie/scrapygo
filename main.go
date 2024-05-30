package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"scrapygo/internal/config"
	"scrapygo/internal/database"
)

type appConfig struct {
	DB     *database.Queries
	Logger *slog.Logger
}

func main() {
	cfg := config.NewConfig()

	config := appConfig{
		DB:     cfg.DB,
		Logger: cfg.Logger,
	}

	mux := config.routes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
