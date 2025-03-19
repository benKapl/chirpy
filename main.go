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
	platform       string
	JWTSecret      string
}

func main() {

	const filepathRoot = "."
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("JWTSecret must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer dbConn.Close()
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		JWTSecret:      JWTSecret,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", middlewareLog(apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.Handle("GET /api/healthz", middlewareLog(http.HandlerFunc(handlerReadiness)))
	mux.Handle("POST /api/users", middlewareLog(http.HandlerFunc(apiCfg.handlerCreateUser)))
	mux.Handle("POST /api/login", middlewareLog(http.HandlerFunc(apiCfg.handlerLogin)))
	mux.Handle("POST /api/refresh", middlewareLog(http.HandlerFunc(apiCfg.handlerRefresh)))
	mux.Handle("POST /api/revoke", middlewareLog(http.HandlerFunc(apiCfg.handlerRevoke)))
	mux.Handle("GET /api/chirps", middlewareLog(http.HandlerFunc(apiCfg.handlerGetChirps)))
	mux.Handle("GET /api/chirps/{id}", middlewareLog(http.HandlerFunc(apiCfg.handlerGetChirp)))
	mux.Handle("POST /api/chirps", middlewareLog(apiCfg.AuthenticateAccessToken(http.HandlerFunc(apiCfg.handlerCreateChirp)))) // used authent middleware
	mux.Handle("GET /admin/metrics", middlewareLog(http.HandlerFunc(apiCfg.handlerMetrics)))
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
