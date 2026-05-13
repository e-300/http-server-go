package main

import (
	"encoding/json"
	"errors"
	"strings"
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
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body     string `json:"body"`
	}

	// get token from request header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT not found in request", err)
		return
	}


	// validate token from header with jwt secret to get the authenticated uuid
	validUserID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "JWT Cant be Authortized", err)
		return
	}
 	
	// unmarshalling r.body stream into decoder then mapping decoder var pointer to params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldnt not Decode request parameters", err)
		return
	}

	// business logic filter out swear words and chirp len
	cleaned, err := validateChirp(params.Body)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	
	// creating a db post object  
	postParams := database.CreatePostParams{
		Body:   cleaned,
		UserID: validUserID,
	}

	// creating creating the post via the post Params object
	post, err := cfg.db.CreatePost(r.Context(), postParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt Create Chirp", err)
		return
	}

	// response with the created chirp db level details 
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        post.ID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Body:      post.Body,
		UserID:    post.UserID,
	})
}


func validateChirp(body string)(string, error){
	const maxChirpLen = 140
	if len(body) > maxChirpLen{
		return "", errors.New("Chirp too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert" : {},
		"fornax"   : {},
	}

	cleaned := getCleanBody(body, badWords)
	return cleaned, nil
}

func getCleanBody(body string, badWords map[string]struct{})(string){
	words := strings.Split(body, " ")
	for i, word := range words{
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok{

			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
