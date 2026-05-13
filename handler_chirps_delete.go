package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/e-300/http-server-go/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldnt find Token", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret) 
	if err != nil{
		respondWithError(w, http.StatusUnauthorized, "Couldnt validate jwt", err)
		return
	}
	
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil{
		respondWithError(w, http.StatusNotFound, "Chirp Could not be Found", err)
		return
	}	


	if dbChirp.UserID != userID{
		respondWithError(w, http.StatusForbidden, "You cant delete this chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpId)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "Chirp could not be deleted", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
