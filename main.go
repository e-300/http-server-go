package main

import (
	"database/sql"
	"github.com/e-300/http-server-go/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

// Stateful handler to track number of requests that have been processed since t-0
// atomic.Int32 -> safely increment and read int val across multiple goroutines or https requests
type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
	polkaSecret    string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	token_string := os.Getenv("JWT_SECRET")
	polka_api_key := os.Getenv("POLKA_SECRET") 

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("DataBase failed to open", err)
	}
	dbQueries := database.New(db)
	const filePathRoot = "."
	const port = "8080"

	// State Object created
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:   token_string,
		polkaSecret:      polka_api_key,
	}

	// NewServeMux -> lookup table matching incoming request -> endpoint -> Handler
	mux := http.NewServeMux()
	// Handle -> expects: Handler Object, http.FilerServer returns one of those handler objects
	// http.FileServer(http.Dir(".") -> This gives us a fileserver that can handle requests
	// Why StripPrefix -> /app/ doesnt exist in our root
	fsHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))
	wrappedHandler := apiCfg.middlewareMetricsInc(fsHandler)
	mux.Handle("/app/", wrappedHandler)


	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)


	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)


	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUserUpdateLogin)

	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	

	mux.HandleFunc("GET /admin/metrics", apiCfg.numOfHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHits)


	// added pointer to server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())

}
