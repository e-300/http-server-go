package main

import (
	"net/http"
	"time"

	"github.com/e-300/http-server-go/internal/auth"

)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){
	//defer r.Body.Close()

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


	if time.Now().After(dbToken.ExpiresAt){
		respondWithError(w, 401, "Token is Expired", err)
		return
	}

	jwt, err := auth.MakeJWT(dbToken.UserID, cfg.token_string, time.Hour)
	if err != nil{
		respondWithError(w, 401, "New jwt could not be issued", err)
		return
	}


	respondWithJSON(w, 200, response{
		Token: jwt,
	})

}