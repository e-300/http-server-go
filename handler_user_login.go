package main

import (
	"encoding/json"
	"net/http"
	"time"

	//"database/sql"

	"github.com/e-300/http-server-go/internal/auth"
	"github.com/e-300/http-server-go/internal/database"
	//"github.com/e-300/http-server-go/internal/database"
)


func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type parameters struct{
		Email string 			  `json:"email"`
		Password  string 	      `json:"password"`
		// Token			string 	  `json:"token"`
		// RefreshToken	string    `json:"refresh_token"`

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

	
	signedAccessToken, err := auth.MakeJWT(user.ID, cfg.token_string, time.Hour)
	if err != nil{
		respondWithError(w, 401, "Access Token Signing Issue", err)
		return
	}

	// creating new refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil{
		respondWithError(w, 401, "Access Token Signing Issue", err)
		return
	}

	_, err = cfg.db.CreateRefresh(r.Context(), database.CreateRefreshParams{
		Token: refreshToken,
		UserID: user.ID, 
	})
	if err != nil{
		respondWithError(w, 401, "Refresh Token could not be inserted into DB", err)
		return
	}


	respondWithJSON(w, 200, response{
		User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: signedAccessToken,
			RefreshToken: refreshToken,
	}})

}