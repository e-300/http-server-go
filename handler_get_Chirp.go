package main

import (
	//"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	
	chirp, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		respondWithError(w, 400, "Couldn't find or parse chirp id", err)
		return
	}

	getChirp, err := cfg.db.GetChirp(r.Context(), chirp)
	if err != nil{
		respondWithError(w, 404, "Chirp not found", err)
	}

	respondWithJSON(w, 200, Chirp{
			ID: getChirp.ID,
			CreatedAt: getChirp.CreatedAt,
			UpdatedAt: getChirp.UpdatedAt,	
			Body: getChirp.Body,     	
			UserID: getChirp.UserID,
		})
	
}