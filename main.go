package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"scrapygo/internal/database"
)

type appConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)

	config := appConfig{
		DB: database.New(db),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/readiness", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, struct{ Status string }{
			Status: "ok",
		})
	})
	mux.HandleFunc("GET /v1/error", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
