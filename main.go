package main


import (
	"net/http"
	"log"
	"sync/atomic"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/Karina-Pogorzelec/Chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db			   *database.Queries
	platform	   string
}

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		db: dbQueries,
		platform: platform,
	}

	serverMux := http.NewServeMux()

	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serverMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)


	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
