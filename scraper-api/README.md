# Frontend RSS Scraper API

Lightweight Go RSS/Atom scraper API for frontend ecosystem updates. It stores feed items in SQLite and exposes them for the Next.js viewer.

## Endpoints

- `GET /api/resources`
- `GET /api/resources?tag=react`
- `GET /api/resources?tag=nextjs&limit=10`
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

## Feed Config

Feeds are configured in `config.json`:

```json
[
  {
    "name": "React",
    "url": "https://react.dev/feed.xml",
    "tag": "react"
  }
]
```
