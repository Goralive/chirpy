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
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
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
	}
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
