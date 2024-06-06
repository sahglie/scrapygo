package scrapygo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"scrapygo/internal/database"
	"scrapygo/internal/validator"
	"time"
)

type createFeedParams struct {
	Name                string `json:"name"`
	Url                 string `json:"url"`
	validator.Validator `json:"-"`
}

type feedParams struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const ErrPgDuplicateFeedUrl = `pq: duplicate key value violates unique constraint "feeds_url_key"`

func (app *application) handlerFeedCreate(w http.ResponseWriter, r *http.Request) {
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

	params.CheckField(validator.NotBlank(params.Name), "name", "can't be blank")
	params.CheckField(validator.NotBlank(params.Url), "url", "can't be blank")
	if !params.Valid() {
		respondWithError(w, http.StatusUnprocessableEntity, params.FirstError())
		return
	}

	feed, err := app.DB.CreateFeed(context.TODO(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    userId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		if err.Error() == ErrPgDuplicateFeedUrl {
			respondWithError(w, http.StatusUnprocessableEntity, "url: already taken")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "failed to create feed")
		return
	}

	feedFollow, err := app.DB.CreateFeedFollow(context.TODO(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    userId,
		FeedID:    feed.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		if err.Error() == ErrPgDuplicateFeedFollow {
			respondWithError(w, http.StatusUnprocessableEntity, "user is already following that feed")
			return
		}

		respondWithError(w, http.StatusInternalServerError, "failed to create feed_follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]any{
		"feed": feedParams{
			ID:        feed.ID,
			Name:      feed.Name,
			Url:       feed.Url,
			UserId:    feed.UserID,
			CreatedAt: feed.CreatedAt,
			UpdatedAt: feed.UpdatedAt,
		},
		"feed_follow": feedFollowParams{
			ID:        feedFollow.ID,
			UserID:    feedFollow.UserID,
			FeedID:    feedFollow.FeedID,
			CreatedAt: feedFollow.CreatedAt,
			UpdatedAt: feedFollow.UpdatedAt,
		},
	})
}

func (app *application) handlerFeedList(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.DB.GetFeeds(r.Context())
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

	respondWithJSON(w, http.StatusOK, map[string][]feedParams{
		"data": feedList,
	})
}
