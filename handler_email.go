package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/e-300/http-server-go/internal/database"
	"github.com/e-300/http-server-go/internal/database/auth"
)


func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	type requestBody struct{
		Email string `json:"email"`
		Password  string    `json:"password"`
	}

	
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
	hash, err := auth.HashPassword(params.Password)
	if err != nil{
		log.Print(err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	})
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

}