package main

import "net/http"

func (cnf *appConfig) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct{ Status string }{
		Status: "ok",
	})
}
func (cnf *appConfig) handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "internal server error")
}
