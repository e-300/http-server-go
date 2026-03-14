package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

// Stateful handler to track number of requests that have been processed since t-0
// atomic.Int32 -> safely increment and read int val across multiple goroutines or https requests
type apiConfig struct{
	fileserverHits atomic.Int32
}

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

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
    w.Write([]byte("Hits reset to 0"))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}


// helper funcs
func respondWithError(w http.ResponseWriter, code int, msg string)

func respondwithJSON(w http.ResponseWriter, code int, payload interface{})

// JSON FUNCTION 1 
func validateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string // key will be msg
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		log.Printf("Error Something went wrong", err)
		w.WriteHeader(500)
		// call respondWithError - pass writer, response code, msg string
		return
	} 

	// payload
	type res struct{
		res bool
	}

	// call respond with json - pass writer, response code , payload interface{}

}

func main(){


	// State Object created
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// NewServeMux -> lookup table matching incoming request -> endpoint -> Handler
	// regisBook is the name i picked instead of typically mux 
	regisBook := http.NewServeMux()
	// Handle -> expects: Handler Object, http.FilerServer returns one of those handler objects
	// http.FileServer(http.Dir(".") -> This gives us a fileserver that can handle requests 
	// Why StripPrefix -> /app/ doesnt exist in our root 
	fsHandler := http.StripPrefix("/app/",http.FileServer(http.Dir(".")))
	wrappedHandler := apiCfg.middlewareMetricsInc(fsHandler)
	regisBook.Handle("/app/", wrappedHandler)


	regisBook.HandleFunc("GET /admin/metrics", apiCfg.numOfHits)

	regisBook.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	regisBook.HandleFunc("GET /api/healthz", handlerReadiness)

	//new endpoint
	regisBook.HandleFunc("POST /api/validate_chirp", validateChirp)
	
	server := http.Server{
		Addr: ":8080",
		Handler: regisBook,
	}

	
	server.ListenAndServe()




	
}


