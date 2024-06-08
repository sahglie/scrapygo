package scrapygo

import (
	"net/http"
)

func (app *Application) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
func (app *Application) handlerError(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "internal server error")
}
