package main

import (
	"net/http"
	"github.com/e-300/http-server-go/internal/auth"
	"time"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){
	type response struct{
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, 401, "Token Missing From Header", err)
		return
	}

	dbToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil{
		respondWithError(w, 401, "Token Not in DB", err)
		return
	}
	if dbToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Token expired", nil)
		return
	}

	if dbToken.RevokedAt.Valid {
		respondWithError(w, 401, "Token revoked", nil)
		return
	}
	
	respondWithJSON(w, 204, http.NoBody)
}