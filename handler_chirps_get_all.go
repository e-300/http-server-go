package main

import (

	"net/http"
)

func (cfg *apiConfig) handlerChirpsAll(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps from database /ncomon mayn", err)
		return
	}

	chirps := []Chirp{}

	for _, chirp := range(dbChirps){
		chirps = append(chirps,Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,	
			Body: chirp.Body,     	
			//UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}