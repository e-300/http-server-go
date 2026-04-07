package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/e-300/http-server-go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Stateful handler to track number of requests that have been processed since t-0
// atomic.Int32 -> safely increment and read int val across multiple goroutines or https requests
type apiConfig struct{
	fileserverHits atomic.Int32
	db 			   *database.Queries 
	platform	   string
	token_string   string
}

func main(){
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	token_string := os.Getenv("TOKEN_STRING")
	
	db, err := sql.Open("postgres", dbURL)
	if err != nil{
		log.Fatal("DataBase failed to open", err)
	}
	dbQueries := database.New(db)
	const filePathRoot = "."
	const port = "8080"


	// State Object created
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db : dbQueries,	
		platform: platform,
		token_string: token_string,
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

	// reset metric hits 
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Validate Chirp length
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	// Create user 
	mux.HandleFunc("POST /api/users" , apiCfg.handlerCreateUser)

	// Create Chirp
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)

	// Get all chirps 
	mux.HandleFunc("GET /api/chirps", apiCfg.getAllChirps)

	// Get single chirp 
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	
	// Login endpoint
	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)

	// added pointer to server
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}



