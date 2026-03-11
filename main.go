package main

import (
	"net/http"
	
)

func main(){


	ptr := http.NewServeMux()

	server := http.Server{
		Addr: ":8080",
		Handler: ptr,
	}


	server.ListenAndServe()
	



	
}


