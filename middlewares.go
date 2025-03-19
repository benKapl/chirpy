package main

import (
	"context"
	"log"
	"net/http"

	"github.com/benKapl/chirpy/internal/auth"
)

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) AuthenticateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
			return
		}
		userID, err := auth.ValidateJWT(token, cfg.JWTSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
			return
		}

		// Store userID in request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
