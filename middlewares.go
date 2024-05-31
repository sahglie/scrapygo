package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
)

//type authorizedUserId uuid.UUID

func (cfg *application) authorizationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := extractApiKey(r.Header.Get("Authorization"))

		user, err := cfg.DB.FindUserByApiKey(r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusUnauthorized, "not authorized")
				return
			}

			msg := fmt.Sprintf("unexpected error: %s\n", err)
			respondWithError(w, http.StatusInternalServerError, msg)
			return
		}

		ctx := context.WithValue(r.Context(), "AuthorizedUserId", user.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
