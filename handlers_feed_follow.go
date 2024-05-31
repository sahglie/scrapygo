package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"scrapygo/internal/authz"
	"scrapygo/internal/database"
	"time"
)

type createFeedFollowParams struct {
	FeedID uuid.UUID `json:"feed_id"`
}

type feedFollowParams struct {
	ID        uuid.UUID `json:"id"`
	FeedID    uuid.UUID `json:"feed_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cfg *application) handlerFeedFollowCreate(w http.ResponseWriter, r *http.Request) {
	userId, err := authz.GetAuthzUser(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	params := createFeedFollowParams{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to decode json")
		return
	}

	feed, err := cfg.DB.CreateFeedFollow(context.TODO(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    userId,
		FeedID:    params.FeedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		fmt.Printf("failed to create feed_follow: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, "failed to create feed_follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollowParams{
		ID:        feed.ID,
		UserID:    feed.UserID,
		FeedID:    feed.FeedID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
	})
}

func (cfg *application) handlerFeedFollowList(w http.ResponseWriter, r *http.Request) {
	userId, err := authz.GetAuthzUser(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	feedFollows, err := cfg.DB.GetUserFeedFollows(r.Context(), userId)
	if err != nil {
		msg := fmt.Sprintf("unexpected error: %s\n", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	feedFollowList := make([]feedFollowParams, len(feedFollows))
	for i, ff := range feedFollows {
		feedFollowList[i] = feedFollowParams{
			ID:        ff.ID,
			UserID:    ff.UserID,
			FeedID:    ff.FeedID,
			CreatedAt: ff.CreatedAt,
			UpdatedAt: ff.UpdatedAt,
		}
	}

	respondWithJSON(w, http.StatusOK, map[string][]feedFollowParams{
		"data": feedFollowList,
	})
}
