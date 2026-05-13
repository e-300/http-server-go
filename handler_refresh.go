package main

import (
	"net/http"
	"time"

	"github.com/e-300/http-server-go/internal/auth"

)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){
	type response struct{
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Couldnt find Token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Token Not in DB", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID, 
		cfg.jwtSecret, 
		time.Hour,
	)

	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldnt issue access Token", err)
		return
	}


	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})

}
