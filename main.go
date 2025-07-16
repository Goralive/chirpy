package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Goralive/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	enviroment     string
	signature      string
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	signature := os.Getenv("SIGNATURE")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	if signature == "" {
		log.Fatal("SIGNATURE must be set")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Something go wrong", err)
	}
	dbQueries := database.New(db)

	const port = "8080"
	mux := http.NewServeMux()

	fileServer := http.FileServer
	path := http.Dir(".")
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		enviroment:     platform,
		signature:      signature,
	}
	fileHandler := http.StripPrefix("/app", fileServer(path))

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(fileHandler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirps)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	log.Printf("Up and running on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
