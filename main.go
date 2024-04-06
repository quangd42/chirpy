package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	filePathRoot = "."
	port         = "8080"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      corsMux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	apiConfig := &apiConfig{
		serverHits: 0,
	}

	mux.Handle("/app/*", apiConfig.middlewareHitInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot)))))
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	}))
	mux.HandleFunc("/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("/reset", apiConfig.handlerReset)

	fmt.Printf("listening on port %s...", port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", cfg.serverHits)
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