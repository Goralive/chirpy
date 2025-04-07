package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()

	fileServer := http.FileServer
	path := http.Dir(".")
	apiCfg := &apiConfig{}
	fileHandler := http.StripPrefix("/app", fileServer(path))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	log.Printf("Up and running on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
