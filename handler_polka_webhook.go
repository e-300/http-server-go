package main

import(
	"net/http"
	"encoding/json"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type UserId struct{ 
		User	string	`json:"user_id"` 
	} 

	type parameters struct{
		Event	string	`json:"event"`
		Data	UserId	`json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 401, "Could not be decoded", err)
		return

	}
	
	if params.Event != "user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return
	}	

	userUUID, err := uuid.Parse(params.Data.User)
	if err != nil{
		respondWithError(w, 404, "Could not parse UUID", err) 
		return
	}

	err = cfg.db.UpgradeUserById(r.Context(), userUUID)
	if err != nil{
		respondWithError(w, 404, "User couldnt not be upgraded", err) 
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

