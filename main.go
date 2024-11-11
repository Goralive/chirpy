package main

import (
	"net/http"
)

func main() {
	const port = "8080"
	mux := http.NewServeMux()

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	server.ListenAndServe()

}
