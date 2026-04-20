package main

import (
	"github.com/e-300/http-server-go/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldnt Find Token", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, 401, "Couldnt Revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
