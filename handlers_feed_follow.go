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

type feedFollowListParams struct {
	Data []feedFollowParams `json:"data"`
}

func (cfg *appConfig) handlerFeedFollowCreate(w http.ResponseWriter, r *http.Request) {
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

//func (cfg *appConfig) handlerFeedList(w http.ResponseWriter, r *http.Request) {
//	feeds, err := cfg.DB.GetFeeds(r.Context())
//	if err != nil {
//		msg := fmt.Sprintf("unexpected error: %s\n", err)
//		respondWithError(w, http.StatusInternalServerError, msg)
//		return
//	}
//
//	feedList := make([]feedParams, len(feeds))
//	for i, f := range feeds {
//		feedList[i] = feedParams{
//			ID:        f.ID,
//			Name:      f.Name,
//			Url:       f.Url,
//			UserId:    f.UserID,
//			CreatedAt: f.CreatedAt,
//			UpdatedAt: f.UpdatedAt,
//		}
//	}
//
//	respondWithJSON(w, http.StatusOK, feedListParams{
//		Data: feedList,
//	})
//}
