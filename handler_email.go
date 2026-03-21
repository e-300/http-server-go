package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type requestBody struct{
		Email string `json:"email"`
	}

	// type responseBody struct {
		// Id string `json:"id"`
		// Created_at time.Time `json:"created_at"`
		// Updated_at time.Time `json:"updated_at"`
		// Email string `json:"email"`
	// }
	
	// Reading raw JSON bytes from request 
	dat, err := io.ReadAll(r.Body)
	if err != nil{
		err := respondWithError(w, 500, "Something went wrong")
		if err != nil{
			log.Println(err)
		}
		return 
	}

	// Raw bytes Mapped into request struct 
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil{
		err := respondWithError(w, 500, "Something went wrong")
		if err != nil{
			log.Println(err)
		}
		return 
	}
	
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil{
		log.Println(err)
		respondWithError(w, 500, "something went wrong broski")
	}
	
	respondWithJSON(w, 201, User{
		ID: user.ID,      
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})



	// respondWithJSON(w, 201, responseBody{
	// 	Id: user.ID.String(),
	// 	Created_at: user.CreatedAt,
	// 	Updated_at: user.UpdatedAt,
	// 	Email: user.Email,
	// })
}