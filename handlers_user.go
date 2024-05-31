package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"scrapygo/internal/database"
	"time"
)

type userParams struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type userListParams struct {
	Data []userParams `json:"data"`
}

func (cfg *application) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	params := userParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json")
		return
	}

	user, err := cfg.DB.CreateUser(context.TODO(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		fmt.Printf("failed to create user: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, userParams{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *application) handlerUserList(w http.ResponseWriter, r *http.Request) {
	apiKey := extractApiKey(r.Header.Get("Authorization"))

	user, err := cfg.DB.FindUserByApiKey(context.TODO(), apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "user not found")
			return
		}

		msg := fmt.Sprintf("unexpected error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	respondWithJSON(w, http.StatusOK, userParams{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
