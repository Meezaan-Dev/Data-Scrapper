# Frontend RSS Resources Hub

Isolated local stack for a Go RSS scraper API and a Next.js resource viewer.

## Run The Sandbox

From this directory:

```bash
make sandbox
```

This starts the containers in the background. Watch logs with:

```bash
make logs
```

If Docker is not installed, install a lightweight container runtime first:

```bash
make install-sandbox-runtime
```

Then open:

- Frontend: http://localhost:3000
- API: http://localhost:8080/api/resources

Trigger a scrape:

```bash
curl -X POST http://localhost:8080/scrape
```

Or:

```bash
make scrape
```

SQLite data is stored in the Docker volume `scrapper_scraper-data`, so it survives container restarts without writing database files into your working tree.

## Stop Or Reset

Stop the stack:

```bash
docker compose down
```

Or:

```bash
make stop
```

Reset the persisted SQLite volume:

```bash
docker compose down -v
```

Or:

```bash
make reset
```

## Services

- `api`: Go, `net/http`, RSS/Atom scraping, SQLite, weekly cron scrape.
- `frontend`: Next.js viewer pointed at `http://api:8080` inside the Docker network.

The app is intentionally local-first and has no auth.
