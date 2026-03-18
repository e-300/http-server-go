package main

import (
	"net/http"
	"sync/atomic"
	"log"
)

// Stateful handler to track number of requests that have been processed since t-0
// atomic.Int32 -> safely increment and read int val across multiple goroutines or https requests
type apiConfig struct{
	fileserverHits atomic.Int32
}

func main(){

	const filePathRoot = "."
	const port = "8080"


	// State Object created
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// NewServeMux -> lookup table matching incoming request -> endpoint -> Handler
	mux := http.NewServeMux()
	// Handle -> expects: Handler Object, http.FilerServer returns one of those handler objects
	// http.FileServer(http.Dir(".") -> This gives us a fileserver that can handle requests 
	// Why StripPrefix -> /app/ doesnt exist in our root 
	fsHandler := http.StripPrefix("/app/",http.FileServer(http.Dir(filePathRoot)))
	wrappedHandler := apiCfg.middlewareMetricsInc(fsHandler)

	mux.Handle("/app/", wrappedHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.numOfHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	// added pointer to server
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}


