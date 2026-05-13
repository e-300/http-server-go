package main

import (
	"encoding/json"
	"net/http"

	"github.com/e-300/http-server-go/internal/auth"
	"github.com/e-300/http-server-go/internal/database"
)

func (cfg *apiConfig) handlerUserUpdateLogin(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email		string	`json:"email"`
		Password	string	`json:"password"`	
	}

	type response struct{
		User
	}

	// Numero Uno always validate user via accesstoken in req header
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "JWT not found", err)
		return
	}
	userID, err := auth.ValidateJWT(accessToken,cfg.jwtSecret) 
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "JWT could not be validated", err)
		return
	}


	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Password couldnt not be hashed", err)
		return
	}

	newCreds := database.UpdateEmailAndPasswordParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	updatedUser, err := cfg.db.UpdateEmailAndPassword(r.Context(), newCreds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{ 
		User:User{
			ID: updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email: updatedUser.Email,
			Red: updatedUser.IsChirpyRed,
		},
	})
}
