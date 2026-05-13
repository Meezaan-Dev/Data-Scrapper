package scraper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigSupportsLegacyFeedArray(t *testing.T) {
	path := writeConfig(t, `[
		{ "name": "React", "url": "https://react.dev/feed.xml", "tag": "react" },
		{ "name": "Next.js", "url": "https://nextjs.org/feed.xml", "tag": "nextjs" },
		{ "name": "React Releases", "url": "https://github.com/facebook/react/releases.atom", "tag": "react" }
	]`)

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if len(config.Feeds) != 3 {
		t.Fatalf("expected 3 feeds, got %d", len(config.Feeds))
	}
	if len(config.Tags) != 2 {
		t.Fatalf("expected 2 deduped tags, got %d", len(config.Tags))
	}
	if config.Tags[0].Tag != "react" || config.Tags[0].Label != "React" {
		t.Fatalf("unexpected first tag: %+v", config.Tags[0])
	}
}

func TestLoadConfigSupportsStructuredConfig(t *testing.T) {
	path := writeConfig(t, `{
		"tags": [{ "tag": "frontend", "label": "Frontend" }],
		"feeds": [{ "name": "React", "url": "https://react.dev/feed.xml", "tag": "frontend" }],
		"links": [{
			"title": "Cursor Changelog",
			"link": "https://cursor.com/changelog",
			"summary": "Official Cursor product updates.",
			"source_name": "Cursor",
			"tag": "frontend",
			"published_at": "2026-05-13T00:00:00Z"
		}]
	}`)

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}

	if len(config.Tags) != 1 || len(config.Feeds) != 1 || len(config.Links) != 1 {
		t.Fatalf("unexpected config counts: tags=%d feeds=%d links=%d", len(config.Tags), len(config.Feeds), len(config.Links))
	}
}

func TestLoadConfigRejectsMissingFeedFields(t *testing.T) {
	path := writeConfig(t, `[{ "name": "React", "tag": "frontend" }]`)

	_, err := LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "must include name, url, and tag") {
		t.Fatalf("expected missing feed field error, got %v", err)
	}
}

func TestLoadConfigRejectsDuplicateTags(t *testing.T) {
	path := writeConfig(t, `{
		"tags": [
			{ "tag": "frontend", "label": "Frontend" },
			{ "tag": "frontend", "label": "Frontend Again" }
		],
		"feeds": [],
		"links": []
	}`)

	_, err := LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "duplicate tag") {
		t.Fatalf("expected duplicate tag error, got %v", err)
	}
}

func TestLoadConfigRejectsUnknownFeedTag(t *testing.T) {
	path := writeConfig(t, `{
		"tags": [{ "tag": "frontend", "label": "Frontend" }],
		"feeds": [{ "name": "AWS", "url": "https://aws.amazon.com/blogs/mobile/feed/", "tag": "aws" }],
		"links": []
	}`)

	_, err := LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "unknown tag") {
		t.Fatalf("expected unknown feed tag error, got %v", err)
	}
}

func TestLoadConfigRejectsIncompleteStaticLink(t *testing.T) {
	path := writeConfig(t, `{
		"tags": [{ "tag": "ai-tools", "label": "AI Tools" }],
		"feeds": [],
		"links": [{
			"title": "Cursor Changelog",
			"link": "https://cursor.com/changelog",
			"source_name": "Cursor",
			"tag": "ai-tools",
			"published_at": "2026-05-13T00:00:00Z"
		}]
	}`)

	_, err := LoadConfig(path)
	if err == nil || !strings.Contains(err.Error(), "link config at index 0") {
		t.Fatalf("expected incomplete link error, got %v", err)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}
