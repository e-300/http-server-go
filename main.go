package main

import (
	"net/http"
	
)

func handlerReadiness(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "text/plain; charset=utf-8",)
	w.WriteHeader(200)
	w.Write([]byte("OK"))

}

func main(){

	// NewServeMux -> lookup table matching incoming request -> endpoint -> Handler
	regisBook := http.NewServeMux()
	// Handle -> expects: Handler Object, http.FilerServer returns one of those handler objects
	// http.FileServer(http.Dir(".") -> This gives us a fileserver that can handle requests 
	// Why StripPrefix -> /app/ doesnt exist in our root 
	regisBook.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	
	// HandleFunc -> expects the name of a fucntion in this case the handlerReadiness 
	// We give HandleFunc our func signature so the server can later call this function with the writer and request
	regisBook.HandleFunc("/healthz", handlerReadiness)
	
	server := http.Server{
		Addr: ":8080",
		Handler: regisBook,
	}

	
	server.ListenAndServe()




	
}


