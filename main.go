package main

import (
	"encoding/json"
	"fmt"
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

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	log.Printf("Up and running on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func healthzHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(response http.ResponseWriter, request *http.Request) {
	html := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	htmlWithCount := fmt.Sprintf(html, cfg.fileserverHits.Load())
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	fmt.Fprint(response, htmlWithCount)
}

func (cfg *apiConfig) resetHandler(response http.ResponseWriter, request *http.Request) {
	cfg.fileserverHits.Store(0)
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Metrics was reseted"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) validateChirpHandler(response http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(500)
		response.Write([]byte("Something went wrong"))
		return
	}

	if len(params.Body) > 140 {
		resp := errorResponse{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(resp)
		if err != nil {
			response.Header().Set("Content-Type", "application/json")
			response.WriteHeader(500)
			response.Write([]byte("Something went wrong"))
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(dat)
		return
	}

	respBody := validResponse{
		Valid: true,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(500)
		response.Write([]byte("Something went wrong"))
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(dat)
}
