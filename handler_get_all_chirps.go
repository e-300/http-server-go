package main

import (

	"net/http"
)

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	
	allChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps from database /ncomon mayn", err)
		return
	}

	responseChirps := []Chirp{}

	for _, chirp := range(allChirps){
		responseChirps = append(responseChirps,Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,	
			Body: chirp.Body,     	
			UserID: chirp.UserID,
		})
	}


	
	respondWithJSON(w, http.StatusOK, responseChirps)
}