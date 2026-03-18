package main

import (
	"net/http"
	"encoding/json"

)

// JSON request helper functions 
func respondWithError(w http.ResponseWriter, code int, msg string) error{
	return respondWithJSON(w, code, map[string]string{"error" : msg})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error{
	// Serlizing payload into json byte slice
	response, err := json.Marshal(payload)
	if err != nil{
		return err
	}
	// Telling Client we are sending back a json response
	w.Header().Set("Content-Type", "application/json")
	// Cors header allowing any origin allowed to recieve this response
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}