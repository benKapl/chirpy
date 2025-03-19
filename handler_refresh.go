package main

import (
	"net/http"
	"time"

	"github.com/benKapl/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Refresh Token in request", err)
		return
	}

	dbToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't match refresh token in database", err)
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired", err)
		return
	}

	if !dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked", err)
		return
	}

	accessToken, err := auth.MakeJWT(dbToken.UserID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't generate access JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
