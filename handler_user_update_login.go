package main

import (
	"net/http"
	"encoding/json"
	"github.com/e-300/http-server-go/internal/auth"
)

func (cfg *apiConfig) handlerUserUpdateLogin(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	
	type parameters struct{
		Token		string	`json:"token"`
		Email		string	`json:"email"`
		Password	string	`json:"password"`	
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
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

	// create sql query to update email and password of the authenticated user
	// hash new password from the request body then update password 
	// errors -> respond with 401
	// everything goes well respond with 200 code and with the newly updated user 
		// omit password in the response



}
