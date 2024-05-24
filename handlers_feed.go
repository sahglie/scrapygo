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

type createFeedParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type feedParams struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type feedListParams struct {
	Data []feedParams `json:"data"`
}

func (cfg *appConfig) handlerFeedCreate(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("AuthorizedUserId").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	params := createFeedParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json")
		return
	}

	feed, err := cfg.DB.CreateFeed(context.TODO(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		fmt.Printf("failed to create feed: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create feed")
		return
	}

	respondWithJSON(w, http.StatusCreated, feedParams{
		ID:        feed.ID,
		Name:      feed.Name,
		Url:       feed.Url,
		UserId:    feed.UserID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
	})
}

func (cfg *appConfig) handlerFeedList(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		msg := fmt.Sprintf("unexpected error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	feedList := make([]feedParams, len(feeds))
	for i, f := range feeds {
		feedList[i] = feedParams{
			ID:        f.ID,
			Name:      f.Name,
			Url:       f.Url,
			UserId:    f.UserID,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		}
	}

	respondWithJSON(w, http.StatusOK, feedListParams{
		Data: feedList,
	})
}
