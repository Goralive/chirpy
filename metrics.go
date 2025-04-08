package main

import (
	"fmt"
	"net/http"
)

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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHandler(response http.ResponseWriter, request *http.Request) {
	if cfg.enviroment != "dev" {
		respondWithError(response, http.StatusForbidden, "Allowed in dev enviroment", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(request.Context())
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Metrics was reseted. Users was deleted"))
}
