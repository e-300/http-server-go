package main

import (
	"encoding/json"
	"io"
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
	UserID    uuid.NullUUID `json:"user_id"`

}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type requestBody struct{
		Msg string `json:"body"`
		User_id string `json:"user_id"`
	}

	// Reading raw JSON bytes from request 
	dat, err := io.ReadAll(r.Body)
	if err != nil{
		respondWithError(w, 500, "Something went wrong broski", err)
		return 
	}
	
	// Raw bytes Mapped into request struct 
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil{
		respondWithError(w, 500, "Something went wrong broski", err)
		return 
	}

	// Checking if user is authenticated 
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, 500, "Something wrong with Bearer", err)
		return 	
	}

	uid, err := auth.ValidateJWT(token, cfg.token_string)
	if err != nil{
		respondWithError(w, 500, "Something went wrong when validating", err)
		return 
	}


	requestMsg := params.Msg
	if len(requestMsg) > 140{
		respondWithError(w, 400, "Chirp is too long dawg", err)
		return 
	}
	res := profaneWords(requestMsg)

	reqUid := params.User_id
	
	parsedUid, err := uuid.Parse(reqUid)
	if err != nil{
		log.Println(err)
		return
	}

	if uid != parsedUid{
		respondWithError(w, 401, "Authentication Could not be established", err)
		return
	}

	cleanUid := uuid.NullUUID{
		UUID: uid,
		Valid: true,
	}



	postParams := database.CreatePostParams{
		Body: res,
		UserID: cleanUid,

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
 		UserID: cleanUid,
	})   
}