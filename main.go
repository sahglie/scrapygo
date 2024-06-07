package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"scrapygo/cmd/scrapygo"
	"scrapygo/internal/config"
	"time"
)

func main() {
	app := scrapygo.Application{
		AppConfig: config.NewConfig(),
	}

	mux := app.Routes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	quit := make(chan bool)
	go app.ScrapeOnInterval(quit, 60*time.Second)

	defer close(quit)

	addr := fmt.Sprintf("%s:%s", host, port)
	app.Logger.Info("starting http server")
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
