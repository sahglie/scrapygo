package scrapygo

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"scrapygo/internal/authz"
	"time"
)

type postParams struct {
	ID          uuid.UUID `json:"id"`
	FeedID      uuid.UUID `json:"feed_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (app *application) handlerPostList(w http.ResponseWriter, r *http.Request) {
	userId, err := authz.GetAuthzUser(r.Context())
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	posts, err := app.DB.GetPostsByUserID(context.TODO(), userId)

	postList := make([]postParams, len(posts))
	for i, p := range posts {
		postList[i] = postParams{
			ID:          p.ID,
			FeedID:      p.FeedID,
			Title:       p.Title,
			Description: p.Description.String,
			Url:         p.Url,
			PublishedAt: p.PublishedAt,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}
	}

	respondWithJSON(w, http.StatusOK, map[string][]postParams{
		"data": postList,
	})
}
