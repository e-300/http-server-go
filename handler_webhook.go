package main

import (
	"encoding/json"
	"net/http"

	"github.com/e-300/http-server-go/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Event	string	`json:"event"`
		Data	struct{
			UserId	uuid.UUID	`json:"user_id"`
		}
	}

	apiKey, err := auth.GetPolkaApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Cant find api key", err)
		return
	}
	if apiKey != cfg.polkaSecret{
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not be decoded", err)
		return

	}
	


	if params.Event != "user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return
	}	

	err = cfg.db.UpgradeUserById(r.Context(), params.Data.UserId)
	if err != nil{
		respondWithError(w, http.StatusInternalServerError, "User couldnt not be upgraded", err) 
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

