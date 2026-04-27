package main

import (
	"github.com/e-300/http-server-go/internal/auth"
	"net/http"
	"github.com/google/uuid"
)
// 1 -> check token in header and make sure user is the author of the chirp
// 1a -> if NOT return 403 status code
// 2 -> delete chirp - so i gotta write a sql query to delete a chirp based on the chirp id from a particular user? i guess check db schema to make sure 
// 2a -> if chirp delete successfully then return 204 status code 
// 2b -> if chirp not found then return a 404 status code 
func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

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


	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil{
		respondWithError(w, 403, "Couldn't find or parse chirp id", err)
		return
	}
	
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
	if err != nil{
		respondWithError(w, 404, "Chirp Could not be Found", err)
		return
	}	


	if dbChirp.UserID != user{
		respondWithError(w, 403, "Userid and Chirp Id Do NOT match", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpId)
	if err != nil{
		respondWithError(w, 403, "Chirp could not be deleted", err)
		return
	}

	w.WriteHeader(204)
}
