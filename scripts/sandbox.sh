#!/usr/bin/env sh
set -eu

if ! command -v docker >/dev/null 2>&1; then
  cat <<'EOF'
Docker is not installed or not available on PATH.

This project sandbox uses containers so the Go API, Next.js app, Node modules,
and SQLite runtime stay isolated from your host environment.

Install a local container runtime, then run this again:

  make install-sandbox-runtime
  make sandbox

If you already installed Docker Desktop, open it once and make sure this works:

  docker --version
EOF
  exit 127
fi

if ! docker compose version >/dev/null 2>&1; then
  cat <<'EOF'
Docker is installed, but the Compose plugin is missing.

Install the sandbox runtime:

  make install-sandbox-runtime

Then run:

  make sandbox
EOF
  exit 127
fi

if ! docker info >/dev/null 2>&1; then
  if command -v colima >/dev/null 2>&1; then
    echo "Docker is installed but not responding. Starting Colima..."
    colima start
  else
    echo "Docker is installed but not responding. Start Docker, then run make sandbox again."
    exit 1
  fi
fi

exec docker compose up --build -d
