package scrapygo

import (
	"net/http"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	authzUser := app.authorizationMiddleware

	mux.HandleFunc("GET /v1/readiness", app.handlerReadiness)
	mux.HandleFunc("GET /v1/error", app.handlerError)
	mux.HandleFunc("POST /v1/users", app.handlerUserCreate)
	mux.HandleFunc("GET /v1/users", app.handlerUserList)
	mux.HandleFunc("POST /v1/feeds", authzUser(app.handlerFeedCreate))
	mux.HandleFunc("GET /v1/feeds", authzUser(app.handlerFeedList))
	mux.HandleFunc("POST /v1/feed_follows", authzUser(app.handlerFeedFollowCreate))
	mux.HandleFunc("GET /v1/feed_follows", authzUser(app.handlerFeedFollowList))
	mux.HandleFunc("GET /v1/posts", authzUser(app.handlerPostList))

	return mux
}
