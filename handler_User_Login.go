package main

import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/e-300/http-server-go/internal/auth"
)


func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type parameters struct{
		Email string 			`json:"email"`
		Password  string 	    `json:"password"`
		Expires_in_seconds *int `json:"expires_in_seconds"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 401, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, 401, "Incorrect email or password", err)
		return

	}

	match , err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match{
		respondWithError(w, 401, "Incorrect email or password", err)
		return		
	}

	signedToken := ""
	if (params.Expires_in_seconds == nil) || (*params.Expires_in_seconds > 3600){
		signedToken, err = auth.MakeJWT(user.ID, cfg.token_string, time.Hour,)
		if err != nil{
			respondWithError(w, 401, "Token Could not be signed", err)
			return
		}

	}else{
		signedToken, err = auth.MakeJWT(user.ID, cfg.token_string, time.Duration(*params.Expires_in_seconds) * time.Second)
		if err != nil{
			respondWithError(w, 401, "Token Could not be signed", err)
			return
		}	
	}

	respondWithJSON(w, 200, response{
		User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: signedToken,
	}})


}