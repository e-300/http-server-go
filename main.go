package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

// Stateful handler to track number of requests that have been processed since t-0
// atomic.Int32 -> safely increment and read int val across multiple goroutines or https requests
type apiConfig struct{
	fileserverHits atomic.Int32
}

// middleware method to increment fileServerHits counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) numOfHits(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w ,"Hits: %v", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits.Store(0)
	w.WriteHeader(200)
    w.Write([]byte("Hits reset to 0"))
}

func main(){

	// State Object created
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// NewServeMux -> lookup table matching incoming request -> endpoint -> Handler
	regisBook := http.NewServeMux()
	// Handle -> expects: Handler Object, http.FilerServer returns one of those handler objects
	// http.FileServer(http.Dir(".") -> This gives us a fileserver that can handle requests 
	// Why StripPrefix -> /app/ doesnt exist in our root 

	fsHandler := http.StripPrefix("/app/",http.FileServer(http.Dir(".")))

	wrappedHandler := apiCfg.middlewareMetricsInc(fsHandler)

	regisBook.Handle("/app/", wrappedHandler)

	regisBook.HandleFunc("GET /metrics", apiCfg.numOfHits)

	regisBook.HandleFunc("/reset", apiCfg.resetHits)
	
	// HandleFunc -> expects the name of a fucntion in this case the handlerReadiness 
	// We give HandleFunc our func signature so the server can later call this function with the writer and request
	regisBook.HandleFunc("GET /healthz", handlerReadiness)
	
	server := http.Server{
		Addr: ":8080",
		Handler: regisBook,
	}

	
	server.ListenAndServe()




	
}


