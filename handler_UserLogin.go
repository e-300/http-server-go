package main

import (
	"encoding/json"
	"net/http"
	"github.com/e-300/http-server-go/internal/auth"
)


func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type parameters struct{
		Email string `json:"email"`
		Password  string    `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 401, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil{
		respondWithError(w, 401, "Incorrect email or password")
		return

	}

	b , err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil{
		respondWithError(w, 401, "Incorrect email or password")
		return		
	}

	if !b{
		respondWithError(w, 401, "Incorrect email or password")
		return
	}


	respondWithJSON(w, 200, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})


}