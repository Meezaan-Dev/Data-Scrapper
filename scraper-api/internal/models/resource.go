package models

import "time"

// Resource is the single API/storage shape for a scraped feed item.
type Resource struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Summary     string    `json:"summary"`
	PublishedAt time.Time `json:"published_at"`
	SourceName  string    `json:"source_name"`
	Tag         string    `json:"tag"`
	CreatedAt   time.Time `json:"created_at"`
}
