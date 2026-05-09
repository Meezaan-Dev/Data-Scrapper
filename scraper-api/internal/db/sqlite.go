package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"frontend-rss-hub/scraper-api/internal/models"
)

const schema = `
CREATE TABLE IF NOT EXISTS resources (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	link TEXT NOT NULL UNIQUE,
	summary TEXT NOT NULL DEFAULT '',
	published_at DATETIME NOT NULL,
	source_name TEXT NOT NULL,
	tag TEXT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_resources_tag_published_at ON resources(tag, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_resources_published_at ON resources(published_at DESC);
`

// Store wraps SQLite access so handlers and scrapers do not need SQL details.
type Store struct {
	db *sql.DB
}

type ResourceQuery struct {
	Tag   string
	Limit int
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	// SQLite handles concurrent readers well, but a single writer is simpler and
	// avoids "database is locked" surprises for this small local app.
	conn.SetMaxOpenConns(1)

	store := &Store{db: conn}
	if err := store.Init(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return store, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Init() error {
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("initialize sqlite schema: %w", err)
	}
	return nil
}

func (s *Store) UpsertResource(resource models.Resource) error {
	if resource.Link == "" {
		return errors.New("resource link is required")
	}

	if resource.PublishedAt.IsZero() {
		resource.PublishedAt = time.Now().UTC()
	}
	if resource.CreatedAt.IsZero() {
		resource.CreatedAt = time.Now().UTC()
	}

	// Link is the natural feed item identity across runs. Updating on conflict
	// keeps titles/summaries fresh without duplicating old resources.
	_, err := s.db.Exec(`
		INSERT INTO resources (title, link, summary, published_at, source_name, tag, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(link) DO UPDATE SET
			title = excluded.title,
			summary = excluded.summary,
			published_at = excluded.published_at,
			source_name = excluded.source_name,
			tag = excluded.tag
	`, resource.Title, resource.Link, resource.Summary, resource.PublishedAt, resource.SourceName, resource.Tag, resource.CreatedAt)
	if err != nil {
		return fmt.Errorf("upsert resource %q: %w", resource.Link, err)
	}

	return nil
}

func (s *Store) ListResources(query ResourceQuery) ([]models.Resource, error) {
	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	// Keep accidental huge reads from making the local UI or API sluggish.
	if limit > 100 {
		limit = 100
	}

	var (
		rows *sql.Rows
		err  error
	)

	if query.Tag != "" {
		rows, err = s.db.Query(`
			SELECT id, title, link, summary, published_at, source_name, tag, created_at
			FROM resources
			WHERE tag = ?
			ORDER BY published_at DESC
			LIMIT ?
		`, query.Tag, limit)
	} else {
		rows, err = s.db.Query(`
			SELECT id, title, link, summary, published_at, source_name, tag, created_at
			FROM resources
			ORDER BY published_at DESC
			LIMIT ?
		`, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("list resources: %w", err)
	}
	defer rows.Close()

	// Return [] instead of null in JSON when there are no rows.
	resources := []models.Resource{}
	for rows.Next() {
		var resource models.Resource
		if err := rows.Scan(
			&resource.ID,
			&resource.Title,
			&resource.Link,
			&resource.Summary,
			&resource.PublishedAt,
			&resource.SourceName,
			&resource.Tag,
			&resource.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan resource: %w", err)
		}
		resources = append(resources, resource)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate resources: %w", err)
	}

	return resources, nil
}
