#!/usr/bin/env sh
set -eu

is_wsl() {
  # WSL should connect to Docker Desktop from Windows instead of installing a
  # separate Docker daemon inside the distro.
  if [ -r /proc/sys/kernel/osrelease ] && grep -qi microsoft /proc/sys/kernel/osrelease; then
    return 0
  fi

  if [ -r /proc/version ] && grep -qi microsoft /proc/version; then
    return 0
  fi

  return 1
}

if is_wsl; then
  cat <<'EOF'
WSL should use Docker Desktop from Windows.

One-time setup:
  1. Install Docker Desktop on Windows:
     https://www.docker.com/products/docker-desktop/
  2. Open Docker Desktop.
  3. Go to Settings > Resources > WSL Integration.
  4. Enable integration for this distro.
  5. Return here and run:
     make sandbox

This keeps the project sandboxed without installing Go, Node, npm, or SQLite in WSL.
EOF
  exit 0
fi

case "$(uname -s)" in
  Darwin)
    if ! command -v brew >/dev/null 2>&1; then
      cat <<'EOF'
Homebrew is required for the macOS Colima bootstrap.

Install Docker Desktop manually, or install Homebrew first:
https://brew.sh
EOF
      exit 127
    fi

    echo "Installing Colima, Docker CLI, Docker Compose, and Buildx..."
    brew install colima docker docker-compose docker-buildx

    DOCKER_CONFIG="${DOCKER_CONFIG:-$HOME/.docker}"
    mkdir -p "$DOCKER_CONFIG/cli-plugins"

    # Homebrew installs Compose/Buildx outside Docker's default plugin folder.
    if [ -x /opt/homebrew/opt/docker-compose/bin/docker-compose ]; then
      ln -sf /opt/homebrew/opt/docker-compose/bin/docker-compose "$DOCKER_CONFIG/cli-plugins/docker-compose"
    elif [ -x /usr/local/opt/docker-compose/bin/docker-compose ]; then
      ln -sf /usr/local/opt/docker-compose/bin/docker-compose "$DOCKER_CONFIG/cli-plugins/docker-compose"
    fi

    if [ -x /opt/homebrew/opt/docker-buildx/bin/docker-buildx ]; then
      ln -sf /opt/homebrew/opt/docker-buildx/bin/docker-buildx "$DOCKER_CONFIG/cli-plugins/docker-buildx"
    elif [ -x /usr/local/opt/docker-buildx/bin/docker-buildx ]; then
      ln -sf /usr/local/opt/docker-buildx/bin/docker-buildx "$DOCKER_CONFIG/cli-plugins/docker-buildx"
    fi

    echo "Starting Colima..."
    colima start --cpu 2 --memory 4 --disk 20
    ;;
  Linux)
    cat <<'EOF'
Linux setup depends on your distro and permissions.

Ubuntu/Debian:
  sudo apt-get update
  sudo apt-get install -y docker.io docker-compose-plugin
  sudo systemctl enable --now docker
  sudo usermod -aG docker "$USER"

Then open a new terminal and run:
  make sandbox

This script does not run sudo package installs automatically so it does not mutate
your machine without a clear system-level approval step.
EOF
    ;;
  *)
    cat <<'EOF'
Install Docker Desktop or a Docker-compatible runtime with Compose, then run:
  make sandbox
EOF
    ;;
esac

echo
echo "Sandbox runtime guidance complete. Start the project with:"
echo "  make sandbox"
