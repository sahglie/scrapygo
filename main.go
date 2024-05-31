package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"scrapygo/internal/config"
)

type application struct {
	*config.AppConfig
}

func main() {
	app := application{
		AppConfig: config.NewConfig(),
	}

	mux := app.routes()

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
