package scrapygo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "Application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("error marshaling json: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, status int, errMsg string) {
	if errMsg == "" {
		errMsg = "something went wrong"
	}

	respondWithJSON(w, status, struct {
		Error string `json:"error"`
	}{
		Error: errMsg,
	})
}

func extractApiKey(authzHeader string) string {
	tokens := strings.Split(authzHeader, " ")
	if len(tokens) != 2 {
		return ""
	}

	return tokens[1]
}
