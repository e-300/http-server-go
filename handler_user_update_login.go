package main

import (
	"encoding/json"
	"net/http"

	"github.com/e-300/http-server-go/internal/auth"
	"github.com/e-300/http-server-go/internal/database"
)

func (cfg *apiConfig) handlerUserUpdateLogin(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	
	type parameters struct{
//		Token		string	`json:"token"`
		Email		string	`json:"email"`
		Password	string	`json:"password"`	
	}

	type response struct{
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, 401, "Couldnt find Token", err)
		return
	}

	user, err := auth.ValidateJWT(accessToken,cfg.token_string) 
	if err != nil{
		respondWithError(w, 401, "jwt not found", err)
		return
	}

	
	
	//Unnessary call can probably get rid of this
	dbUser, err := cfg.db.GetUserFromId(r.Context(), user)
	if err != nil{
		respondWithError(w, 401, "user could not be found in db", err)
		return
	}
	
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil{
		respondWithError(w, 401, "Password couldt not be hashed", err)
		return
	}

	newCreds := database.UpdateEmailAndPasswordParams{
		ID: dbUser.ID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	}

	updatedUser, err := cfg.db.UpdateEmailAndPassword(r.Context(), newCreds)
	if err != nil {
		respondWithError(w, 401, "Couldnt not update to New Creds", err)
		return
	}

	respondWithJSON(w, 200, response{ 
		User:User{
			ID: updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email: updatedUser.Email,
			},
	})
}
