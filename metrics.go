package main

import (
	"net/http"
	"fmt"
)

// MiddlewareMetrics is a function that returns a Handler obj
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)

		next.ServeHTTP(w,r)
	})
}

// numOfHits, resetHits, handlerReadiness are all "handler functions" 
// bc func signature matches func(w http.ResponseWriter, r *http.Request)
// these will be registerd with the mux with the http.HandlerFunc method
const adminTemplate = "Welcome, Chirpy Admin \nChirpy has been visited %d times!"
func (cfg *apiConfig) numOfHits(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//hits := fmt.Sprintf(adminTemplate,cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w ,adminTemplate,cfg.fileserverHits.Load())
}