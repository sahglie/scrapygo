package main

import (
	"context"
	"encoding/json"
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

func (cnf *appConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	params := userParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json")
		return
	}

	user, err := cnf.DB.CreateUser(context.TODO(), database.CreateUserParams{
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

func (cnf *appConfig) handlerUserList(w http.ResponseWriter, r *http.Request) {
	users, err := cnf.DB.ListUsers(context.TODO())
	if err != nil {
		fmt.Printf("failed to retrieve users: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve users")
		return
	}

	userList := userListParams{Data: make([]userParams, len(users))}
	for i, u := range users {
		userList.Data[i] = userParams{
			ID:        u.ID,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}
	}

	respondWithJSON(w, http.StatusCreated, userList)
}
