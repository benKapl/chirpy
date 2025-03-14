package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/benKapl/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer dbConn.Close()

	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.Handle("GET /api/healthz", middlewareLog(http.HandlerFunc(handlerReadiness)))
	mux.Handle("POST /api/validate_chirp", middlewareLog(http.HandlerFunc(handlerValidateChirp)))
	mux.Handle("GET /admin/metrics", middlewareLog(http.HandlerFunc(apiCfg.handlerMetrics)))
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
