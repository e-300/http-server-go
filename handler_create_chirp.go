package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/e-300/http-server-go/internal/auth"
	"github.com/e-300/http-server-go/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string	`json:"body"`
	UserID    uuid.UUID `json:"user_id"`

}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct{
		Msg string `json:"body"`
		User_id string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Checking if user is authenticated 
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, 500, "Something wrong with Bearer", err)
		fmt.Fprintln(w,token)
		fmt.Fprintln(w,params.User_id)
		return 	
	}
	
	log.Printf("DEBUG raw token: %q", token) 

	validatedUid, err := auth.ValidateJWT(token, cfg.token_string)
	if err != nil{
		respondWithError(w, 401, "Something went wrong when validating", err)
		return 
	}

	requestMsg := params.Msg
	if len(requestMsg) > 140{
		respondWithError(w, 401, "Chirp is too long dawg", err)
		return 
	}
	cleanedMsg := profaneWords(requestMsg)


	postParams := database.CreatePostParams{
		Body: cleanedMsg,
		UserID: validatedUid,

	}

	post , err := cfg.db.CreatePost(r.Context(), postParams)
	if err != nil{
		log.Println(err)
		return
	}
	
	respondWithJSON(w, 201, Chirp{
		ID: post.ID,      	
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
 		Body: post.Body,
 		UserID: post.UserID,
	})   
}