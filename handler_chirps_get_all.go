package main

import (
	"net/http"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsAll(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		authorID = "Default"
	}

	// TO DOOO
	// gotta write a sql query to sort chirps based on author id
	//	by created at in ascedning order
	// call that function if authorID != Default
	// then return the sorted list of chirps
	if authorID != "Default"{
		u, err := uuid.Parse(authorID)
		if err != nil{
			respondWithError(w, http.StatusInternalServerError, "uuid could not be parsed", err)
			return
		}
		dbPosts, err := cfg.db.GetChirpsForAuthor(r.Context(), u)
		if err != nil{
			respondWithError(w, http.StatusInternalServerError, "Could not find database of posts for user", err)
			return
		}
		chirps := []Chirp{}
		for _, chirp := range dbPosts{
			chirps = append(chirps, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				//UserID: chirp.UserID,
			})
		}
		respondWithJSON(w,http.StatusOK, chirps)
		return
	}	


	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get all chirps from database /ncomon mayn", err)
		return
	}

	chirps := []Chirp{}

	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			//UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
