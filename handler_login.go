package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/benKapl/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	const defaultExpirationSeconds = 60 * 60

	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds,omitempty"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	var expirationSeconds int
	// If client doesnt specify expiration
	if params.ExpiresInSeconds == nil {
		expirationSeconds = defaultExpirationSeconds
	} else {
		expirationSeconds = *params.ExpiresInSeconds
		// if client has duration above 1 hour
		if expirationSeconds > defaultExpirationSeconds {
			expirationSeconds = defaultExpirationSeconds
		}
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JWTSecret, time.Duration(expirationSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error creating JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	})
}
