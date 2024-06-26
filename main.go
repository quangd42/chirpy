package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ammon134/chirpy/internal/database"
	"github.com/joho/godotenv"
)

const (
	filePathRoot = "."
	port         = "8080"
	dbPath       = "database.json"
)

type apiConfig struct {
	db           *database.DB
	jwtSecret    string
	polka_apikey string
	serverHits   int
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	apiConfig := &apiConfig{
		db:           db,
		jwtSecret:    os.Getenv("JWT_SECRET"),
		polka_apikey: os.Getenv("POLKA_APIKEY"),
		serverHits:   0,
	}

	mux.Handle("/app/*", apiConfig.middlewareHitInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))

	mux.Handle("GET /api/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}))
	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("GET /api/reset", apiConfig.handlerReset)

	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiConfig.handlerDeleteChirp)

	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiConfig.handlerUpdateUser)

	mux.HandleFunc("POST /api/login", apiConfig.handlerLogin)

	mux.HandleFunc("POST /api/revoke", apiConfig.handlerRevokeToken)
	mux.HandleFunc("POST /api/refresh", apiConfig.handlerRefreshToken)

	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlerWebhookUpgradeUser)

	fmt.Printf("listening on port %s...\n", port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	html := `
  <html>
    <body>
      <h1>Welcome, Chirpy Admin</h1>
      <p>Chirpy has been visited %d times!</p>
    </body>
  </html>
  `
	fmt.Fprintf(w, html, cfg.serverHits)
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.serverHits = 0
	fmt.Fprint(w, "Hits reset to 0")
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
