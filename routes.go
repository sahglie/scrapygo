package main

import "net/http"

func (cfg *application) routes() http.Handler {
	mux := http.NewServeMux()

	authzUser := cfg.authorizationMiddleware

	mux.HandleFunc("GET /v1/readiness", cfg.handlerReadiness)
	mux.HandleFunc("GET /v1/error", cfg.handlerError)
	mux.HandleFunc("POST /v1/users", cfg.handlerUserCreate)
	mux.HandleFunc("GET /v1/users", cfg.handlerUserList)
	mux.HandleFunc("POST /v1/feeds", authzUser(cfg.handlerFeedCreate))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerFeedList)
	mux.HandleFunc("POST /v1/feed_follows", authzUser(cfg.handlerFeedFollowCreate))
	mux.HandleFunc("GET /v1/feed_follows", authzUser(cfg.handlerFeedFollowList))

	return mux
}
