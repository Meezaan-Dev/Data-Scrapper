package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"frontend-rss-hub/scraper-api/internal/db"
	"frontend-rss-hub/scraper-api/internal/scraper"
)

type Handler struct {
	store      *db.Store
	configPath string
}

func New(store *db.Store, configPath string) *Handler {
	return &Handler{store: store, configPath: configPath}
}

func (h *Handler) ListResources(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, err := parseLimit(query.Get("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	resources, err := h.store.ListResources(db.ResourceQuery{
		Tag:   query.Get("tag"),
		Limit: limit,
	})
	if err != nil {
		log.Printf("list resources failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to fetch resources")
		return
	}

	writeJSON(w, http.StatusOK, resources)
}

func (h *Handler) Scrape(w http.ResponseWriter, r *http.Request) {
	feeds, err := scraper.LoadFeeds(h.configPath)
	if err != nil {
		log.Printf("load feeds failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to load feed config")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Minute)
	defer cancel()

	count, err := scraper.New(h.store).Scrape(ctx, feeds)
	if err != nil {
		log.Printf("scrape failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to scrape feeds")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":          "ok",
		"processed_items": count,
	})
}

func parseLimit(value string) (int, error) {
	if value == "" {
		return 20, nil
	}

	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 {
		return 0, strconv.ErrSyntax
	}
	if limit > 100 {
		return 100, nil
	}
	return limit, nil
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
