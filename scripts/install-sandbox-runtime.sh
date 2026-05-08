#!/usr/bin/env sh
set -eu

if ! command -v brew >/dev/null 2>&1; then
  cat <<'EOF'
Homebrew is required for this bootstrap script.

Install Docker Desktop manually, or install Homebrew first:
https://brew.sh
EOF
  exit 127
fi

echo "Installing Colima, Docker CLI, Docker Compose, and Buildx..."
brew install colima docker docker-compose docker-buildx

DOCKER_CONFIG="${DOCKER_CONFIG:-$HOME/.docker}"
mkdir -p "$DOCKER_CONFIG/cli-plugins"

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

echo
echo "Sandbox runtime is ready. Start the project with:"
echo "  make sandbox"
