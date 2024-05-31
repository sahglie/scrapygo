package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	authzUser := app.authorizationMiddleware

	mux.HandleFunc("GET /v1/readiness", app.handlerReadiness)
	mux.HandleFunc("GET /v1/error", app.handlerError)
	mux.HandleFunc("POST /v1/users", app.handlerUserCreate)
	mux.HandleFunc("GET /v1/users", app.handlerUserList)
	mux.HandleFunc("POST /v1/feeds", authzUser(app.handlerFeedCreate))
	mux.HandleFunc("GET /v1/feeds", app.handlerFeedList)
	mux.HandleFunc("POST /v1/feed_follows", authzUser(app.handlerFeedFollowCreate))
	mux.HandleFunc("GET /v1/feed_follows", authzUser(app.handlerFeedFollowList))

	return mux
}
