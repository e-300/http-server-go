package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		respondWithError(w, 400, "Couldn't find or parse chirp id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil{
		respondWithError(w, 404, "Chirp not found", err)
		return
	}

	respondWithJSON(w, 200, Chirp{
			ID: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,	
			Body: dbChirp.Body,     	
			//UserID: dbChirp.UserID,
		})
	
}