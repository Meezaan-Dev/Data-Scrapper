# Frontend RSS Resources Hub

Isolated local stack for a Go RSS scraper API and a Next.js resource viewer.

The project is designed so a teammate can clone it and run one command on macOS, Windows, or Linux. The only host requirement is a Docker-compatible container runtime with Docker Compose, such as Docker Desktop.

## Run The Sandbox

From this directory:

```bash
docker compose up --build
```

That command builds and runs both services, creates the SQLite volume, and performs an initial scrape the first time the database is empty.

Open:

- Frontend: http://localhost:3000
- API: http://localhost:8080/api/resources

To run in the background instead:

```bash
docker compose up --build -d
```

Watch logs:

```bash
docker compose logs -f
```

Trigger a manual scrape:

```bash
curl -X POST http://localhost:8080/scrape
```

SQLite data is stored in the Docker volume `frontend-rss-hub_scraper-data`, so it survives container restarts without writing database files into your working tree.

## One-Command Helpers

On macOS/Linux with `make` available:

```bash
make sandbox
```

On Windows PowerShell, use Docker Compose directly:

```powershell
docker compose up --build
```

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

## Host Requirement

This repository intentionally does not require Go, Node, npm, SQLite, or cron on the host. Those live inside containers.

A container runtime itself still has to exist on the host because no repository can start a sandbox VM/container engine from nothing on every operating system. Install one of:

- Docker Desktop for macOS, Windows, or Linux
- Colima plus Docker CLI on macOS
- Docker Engine plus Docker Compose on Linux

## Services

- `api`: Go, `net/http`, RSS/Atom scraping, SQLite, weekly cron scrape.
- `frontend`: Next.js viewer pointed at `http://api:8080` inside the Docker network.

The app is intentionally local-first and has no auth.
