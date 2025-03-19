package main

import (
	"net/http"

	"github.com/benKapl/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Refresh Token in request", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
