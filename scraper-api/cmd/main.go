package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	"frontend-rss-hub/scraper-api/internal/db"
	"frontend-rss-hub/scraper-api/internal/handlers"
	"frontend-rss-hub/scraper-api/internal/scraper"
)

const (
	defaultAddr       = ":8080"
	defaultConfigPath = "config.json"
	defaultDBPath     = "data/resources.db"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	configPath := envOrDefault("CONFIG_PATH", defaultConfigPath)
	dbPath := envOrDefault("DB_PATH", defaultDBPath)

	// The same binary runs both the long-lived API and one-off scrape jobs.
	// This keeps Docker cron simple: it can call `scraper-api scrape`.
	switch os.Args[1] {
	case "serve":
		if err := serve(dbPath, configPath); err != nil {
			log.Fatal(err)
		}
	case "scrape":
		if err := runScrape(dbPath, configPath); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
		os.Exit(1)
	}
}

func serve(dbPath, configPath string) error {
	// Opening the store also creates the schema if this is a fresh SQLite file.
	store, err := db.Open(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	handler := handlers.New(store, configPath)
	router := mux.NewRouter()
	router.HandleFunc("/api/resources", handler.ListResources).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/tags", handler.ListTags).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/scrape", handler.Scrape).Methods(http.MethodPost, http.MethodOptions)

	addr := envOrDefault("ADDR", defaultAddr)
	server := &http.Server{
		Addr:              addr,
		Handler:           cors(router),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("scraper API listening on %s", addr)
	return server.ListenAndServe()
}

func runScrape(dbPath, configPath string) error {
	// The CLI scrape path uses the same scraper/store code as POST /scrape.
	store, err := db.Open(dbPath)
	if err != nil {
		return err
	}
	defer store.Close()

	config, err := scraper.LoadConfig(configPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	count, err := scraper.New(store).ScrapeConfig(ctx, config)
	if err != nil {
		return err
	}

	log.Printf("processed %d feed items", count)
	return nil
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Local-first app: allow the Next.js viewer to call the API from localhost.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func usage() {
	fmt.Println("usage: scraper-api [serve|scrape]")
}
