# Frontend RSS Scraper API

Lightweight Go RSS/Atom scraper API for frontend ecosystem updates. It stores feed items in SQLite and exposes them for the Next.js viewer.

## Endpoints

- `GET /api/resources`
- `GET /api/resources?tag=frontend`
- `GET /api/resources?tag=llm&limit=10`
- `GET /api/tags`
- `POST /scrape`

## Local Development

```bash
go mod download
go run ./cmd/main.go serve
```

In another terminal:

```bash
curl -X POST http://localhost:8080/scrape
curl http://localhost:8080/api/resources
curl http://localhost:8080/api/tags
```

## CLI

```bash
go run ./cmd/main.go serve
go run ./cmd/main.go scrape
```

Optional environment variables:

- `ADDR` defaults to `:8080`
- `CONFIG_PATH` defaults to `config.json`
- `DB_PATH` defaults to `data/resources.db`

## Docker

```bash
docker build -t frontend-scraper .
docker run -p 8080:8080 -v ./data:/root/data frontend-scraper
```

The container installs a cron job that runs every Monday at 09:00 and stores SQLite data at `/root/data/resources.db`.

## Resource Config

Resources are configured in `config.json`. The current format supports ordered tags, RSS/Atom feeds, and curated static official links:

```json
{
  "tags": [{ "tag": "frontend", "label": "Frontend" }],
  "feeds": [{ "name": "React", "url": "https://react.dev/feed.xml", "tag": "frontend" }],
  "links": [
    {
      "title": "Cursor Changelog",
      "link": "https://cursor.com/changelog",
      "summary": "Official Cursor product updates and release notes.",
      "source_name": "Cursor",
      "tag": "ai-tools",
      "published_at": "2026-05-13T00:00:00Z"
    }
  ]
}
```

The legacy array-only feed format is still accepted for compatibility.
