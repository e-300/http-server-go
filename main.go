package main

import (
	"encoding/json"
	"fmt"
	"io"
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
func respondWithError(w http.ResponseWriter, code int, msg string) error{
	return respondWithJSON(w, code, map[string]string{"error" : msg})
}

// Generic Helper
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


// JSON FUNCTION 1 
func validateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct{
		Msg string `json:"body"`
	}

	type responseBody struct {
		Valid bool `json:"valid"`
	}
	

	dat, err := io.ReadAll(r.Body)
	if err != nil{
		err := respondWithError(w, 500, "Something went wrong")
		if err != nil{
			log.Println(err)
		}
		return 
	}
	
	params := requestBody{}
	err = json.Unmarshal(dat, &params)
	if err != nil{
		err := respondWithError(w, 500, "Something went wrong")
		if err != nil{
			log.Println(err)
		}
		return 
	}


	if len(params.Msg) > 140{
		err := respondWithError(w, 400, "Chirp is too long")
		if err != nil{
			log.Println(err)
		}
		return 
	}

	respondWithJSON(w, 200, responseBody{
		Valid: true,
	})

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


