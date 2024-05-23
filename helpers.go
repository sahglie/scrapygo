package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type userJSON struct {
	Name string `json:"name"`
}

type errorJSON struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

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

	err := errorJSON{Error: errMsg}
	respondWithJSON(w, status, err)
}
