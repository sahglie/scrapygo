package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"scrapygo/internal/database"
	"time"
)

type userParams struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (cnf *appConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	params := userParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json")
		return
	}

	user, err := cnf.DB.CreateUser(context.TODO(), database.CreateUserParams{
		ID:        1,
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
		ID:        string(user.ID),
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
