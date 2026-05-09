#!/usr/bin/env sh
set -eu

wait_for_docker() {
  seconds="${1:-90}"
  i=0

  # Docker Desktop/Colima can take a while to expose the daemon after launching.
  while [ "$i" -lt "$seconds" ]; do
    if docker info >/dev/null 2>&1; then
      return 0
    fi
    sleep 2
    i=$((i + 2))
  done

  return 1
}

is_wsl() {
  # WSL reports Microsoft in kernel metadata; use both files for broad distro support.
  if [ -r /proc/sys/kernel/osrelease ] && grep -qi microsoft /proc/sys/kernel/osrelease; then
    return 0
  fi

  if [ -r /proc/version ] && grep -qi microsoft /proc/version; then
    return 0
  fi

  return 1
}

print_missing_docker() {
  if is_wsl; then
    cat <<'EOF'
Docker is not available in this WSL distro.

Fastest setup:
  1. Install Docker Desktop on Windows.
  2. Open Docker Desktop.
  3. Go to Settings > Resources > WSL Integration.
  4. Enable integration for this distro.
  5. Run: ./sandbox

After that, this repo needs no Go, Node, npm, SQLite, or local app setup.
EOF
    return
  fi

  case "$(uname -s)" in
    Darwin)
      cat <<'EOF'
Docker is not installed or not available on PATH.

Install one container runtime:
  - Docker Desktop for Mac, then run: ./sandbox
  - or Colima/Docker CLI with: ./scripts/install-sandbox-runtime.sh
EOF
      ;;
    Linux)
      cat <<'EOF'
Docker is not installed or not available on PATH.

Install Docker Engine with the Compose plugin, then run:
  ./sandbox

Ubuntu/Debian quick path:
  sudo apt-get update
  sudo apt-get install -y docker.io docker-compose-plugin
  sudo systemctl enable --now docker
EOF
      ;;
    *)
      cat <<'EOF'
Docker is not installed or not available on PATH.

Install Docker Desktop or a Docker-compatible runtime with Compose, then run:
  ./sandbox
EOF
      ;;
  esac
}

try_start_docker() {
  if docker info >/dev/null 2>&1; then
    return 0
  fi

  if is_wsl && command -v powershell.exe >/dev/null 2>&1; then
    # In WSL, Docker usually runs in Windows via Docker Desktop, not inside Linux.
    echo "Detected WSL. Trying to start Docker Desktop on Windows..."
    powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "\$p='C:\Program Files\Docker\Docker\Docker Desktop.exe'; if (Test-Path \$p) { Start-Process -FilePath \$p; exit 0 } else { exit 1 }" >/dev/null 2>&1 || true

    if wait_for_docker 120; then
      return 0
    fi

    cat <<'EOF'
Docker Desktop did not become available to WSL.

Open Docker Desktop on Windows, then check:
  Settings > Resources > WSL Integration

Enable this distro, apply/restart, then run:
  ./sandbox
EOF
    return 1
  fi

  if [ "$(uname -s)" = "Darwin" ]; then
    if command -v colima >/dev/null 2>&1; then
      echo "Starting Colima..."
      colima start
      return 0
    fi

    if command -v open >/dev/null 2>&1; then
      echo "Trying to start Docker Desktop..."
      open -ga Docker >/dev/null 2>&1 || true
      wait_for_docker 120 && return 0
    fi
  fi

  if [ "$(uname -s)" = "Linux" ] && command -v systemctl >/dev/null 2>&1; then
    # Native Linux can often start the daemon directly. sudo may prompt the user.
    echo "Docker is installed but not running. Trying to start Docker..."
    if systemctl start docker >/dev/null 2>&1 || sudo systemctl start docker; then
      wait_for_docker 30 && return 0
    fi

    cat <<'EOF'
Docker is installed but the daemon is not running.

Start it, then rerun this command:
  sudo systemctl enable --now docker
  ./sandbox
EOF
    return 1
  fi

  cat <<'EOF'
Docker is installed but the daemon is not responding.

Start your Docker runtime, then run:
  ./sandbox
EOF
  return 1
}

if ! command -v docker >/dev/null 2>&1; then
  print_missing_docker
  exit 127
fi

if ! docker compose version >/dev/null 2>&1; then
  cat <<'EOF'
Docker is installed, but the Compose plugin is missing.

Install Docker Compose, then run:
  ./sandbox
EOF
  exit 127
fi

try_start_docker

# Compose is the sandbox boundary: Go, Node, SQLite, cron, and app processes all
# run in containers after this point.
docker compose up --build -d

cat <<'EOF'

Sandbox is starting.

Open:
  http://localhost:3000/resources

Useful commands:
  make logs
  make scrape
  make stop
EOF
