package main

import (
	"encoding/json"
	"net/http"
	"io"
	"strings"
)

// Profane Word Func
func profaneWords(r string) string{
	profaneWords := map[string]string{
		"kerfuffle": "****",
		"sharbert": "****",
		"fornax": "****",
	
	}
	res := ""
	wordList := strings.Split(r, " ")
	for _, word := range wordList{
		elem, ok := profaneWords[string(strings.ToLower(word))]
		if ok{
			res += elem
		}else{
			res += string(word)
		}
		res += " "
	}
	res = strings.TrimRight(res, " ")
	return res
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct{
		Msg string `json:"body"`
	}

	type responseBody struct {
		Valid bool `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}
	
	// Reading raw JSON bytes from request 
	dat, err := io.ReadAll(r.Body)
	if err != nil{
		respondWithError(w, 500, "Something went wrong", err)
		return 
	}
	
	// Raw bytes Mapped into request struct 
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil{
		respondWithError(w, 500, "Something went wrong", err)
		return 
	}

	// Business Logic

	requestMsg := params.Msg
	if len(requestMsg) > 140{
		respondWithError(w, 400, "Chirp is too long", err)
		return 
	}
	res := profaneWords(requestMsg)
	respondWithJSON(w, 200, responseBody{
		Valid: true,
		CleanedBody: res,
	},)

}