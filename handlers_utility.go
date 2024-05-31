package main

import "net/http"

func (app *application) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
func (app *application) handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "internal server error")
}
