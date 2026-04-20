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
		respondWithError(w, 401, "Couldnt find Token", err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil{
		respondWithError(w, 401, "Token Not in DB", err)
		return
	}

//	// token is expired???
//	if time.Now().After(dbToken.ExpiresAt){
//		respondWithError(w, 401, "Token is Expired", err)
//		return
//	}
//
//	// token is revoked??
//	if dbToken.RevokedAt.Valid == true{
//		respondWithError(w, 401, "Token is Revoked", err)
//		return
//	}
	// issue new jwt
	accessToken, err := auth.MakeJWT(
		user.ID, 
		cfg.token_string, 
		time.Hour,
	)

	if err != nil{
		respondWithError(w, 401, "Couldnt issue access Token", err)
		return
	}


	respondWithJSON(w, 200, response{
		Token: accessToken,
	})

}
