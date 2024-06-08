package scrapygo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"net/http"
	"scrapygo/internal/database"
	"scrapygo/internal/validator"
	"time"
)

type userParams struct {
	ID                  uuid.UUID `json:"id,omitempty"`
	Name                string    `json:"name"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
	validator.Validator `json:"-"`
}

const ErrPgDuplicateUserName = `pq: duplicate key value violates unique constraint "unique_name_idx"`

func (app *Application) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	params := userParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json")
		return
	}

	params.CheckField(validator.NotBlank(params.Name), "name", "can't be blank")

	if !params.Valid() {
		respondWithError(w, http.StatusUnprocessableEntity, params.FirstError())
		return
	}

	user, err := app.DB.CreateUser(context.TODO(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		if err.Error() == ErrPgDuplicateUserName {
			respondWithError(w, http.StatusUnprocessableEntity, "name: already taken")
			return
		}

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

func (app *Application) handlerUserList(w http.ResponseWriter, r *http.Request) {
	apiKey := extractApiKey(r.Header.Get("Authorization"))

	user, err := app.DB.FindUserByApiKey(context.TODO(), apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "not authorized")
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
