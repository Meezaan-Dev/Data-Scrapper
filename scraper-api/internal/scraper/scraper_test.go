package scraper

import (
	"path/filepath"
	"testing"
	"time"

	"frontend-rss-hub/scraper-api/internal/db"
)

func TestUpsertLinksDoesNotDuplicateResources(t *testing.T) {
	store, err := db.Open(filepath.Join(t.TempDir(), "resources.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer store.Close()

	link := LinkConfig{
		Title:       "Cursor Changelog",
		Link:        "https://cursor.com/changelog",
		Summary:     "Official Cursor product updates.",
		SourceName:  "Cursor",
		Tag:         "ai-tools",
		PublishedAt: time.Date(2026, 5, 13, 0, 0, 0, 0, time.UTC),
	}

	scraper := New(store)
	for i := 0; i < 2; i++ {
		count, err := scraper.UpsertLinks([]LinkConfig{link})
		if err != nil {
			t.Fatalf("upsert links pass %d: %v", i+1, err)
		}
		if count != 1 {
			t.Fatalf("expected 1 processed link on pass %d, got %d", i+1, count)
		}
	}

	resources, err := store.ListResources(db.ResourceQuery{Tag: "ai-tools", Limit: 100})
	if err != nil {
		t.Fatalf("list resources: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 stored resource after duplicate upserts, got %d", len(resources))
	}
	if resources[0].Link != link.Link {
		t.Fatalf("unexpected stored resource link: %s", resources[0].Link)
	}
}
