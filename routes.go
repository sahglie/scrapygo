package main

import "net/http"

func (cfg *appConfig) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/readiness", cfg.handlerReadiness)
	mux.HandleFunc("GET /v1/error", cfg.handlerError)
	mux.HandleFunc("POST /v1/users", cfg.handlerUserCreate)
	mux.HandleFunc("GET /v1/users", cfg.handlerUserList)
	mux.HandleFunc("POST /v1/feeds", cfg.authorizationMiddleware(cfg.handlerFeedCreate))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerFeedList)
	mux.HandleFunc("POST /v1/feed_follows", cfg.authorizationMiddleware(cfg.handlerFeedFollowCreate))

	return mux
}
