# Frontend RSS Resources Hub

Isolated local stack for a Go RSS scraper API and a Next.js resource viewer.

The project is designed so a teammate can clone it and run one command on macOS, Windows, or Linux. The only host requirement is a Docker-compatible container runtime with Docker Compose, such as Docker Desktop.

## Run The Sandbox

From this directory on macOS, Linux, or WSL:

```bash
./sandbox
```

That command checks Docker, gives platform-specific setup help when Docker is not ready, builds and runs both services, creates the SQLite volume, and performs an initial scrape the first time the database is empty.

From Windows PowerShell:

```powershell
.\sandbox.ps1
```

From Windows Command Prompt:

```bat
sandbox.cmd
```

Open:

- Frontend: http://localhost:3000
- API: http://localhost:8080/api/resources

If you prefer raw Docker Compose:

```bash
docker compose up --build
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

## First-Time Runtime Setup

This repository intentionally does not require Go, Node, npm, SQLite, or cron on the host. Those live inside containers.

A container runtime itself still has to exist on the host because no repository can start a VM/container engine from nothing on every operating system.

Recommended setup by environment:

- macOS: Docker Desktop, or run `./scripts/install-sandbox-runtime.sh` to install Colima with Homebrew.
- Windows PowerShell: Docker Desktop, then `.\sandbox.ps1`.
- WSL: Docker Desktop on Windows with WSL integration enabled, then `./sandbox`.
- Linux: Docker Engine plus the Docker Compose plugin.

### WSL

For WSL, use Docker Desktop from Windows:

1. Install and open Docker Desktop on Windows.
2. Go to Settings > Resources > WSL Integration.
3. Enable integration for your distro.
4. In WSL, run `./sandbox`.

If Docker Desktop is installed but closed, `./sandbox` will try to start it from WSL and wait for Docker to become available.

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
