#!/usr/bin/env bash

set -Eeuo pipefail

script_dir="$(CDPATH='' cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
project_dir="$(dirname -- "$script_dir")"
compose=(docker compose --file "$project_dir/docker-compose.yaml")

usage() {
  printf '%s\n' \
    "Usage: ./scripts/setup.sh [--help]" \
    "" \
    "Build and start the complete library stack with Docker Compose."
}

case "${1:-}" in
  "") ;;
  --help|-h)
    usage
    exit 0
    ;;
  *)
    usage >&2
    exit 2
    ;;
esac

if (( $# > 1 )); then
  usage >&2
  exit 2
fi

if ! command -v docker >/dev/null 2>&1; then
  printf '%s\n' "Docker is required. Install Docker Desktop or Docker Engine, then retry." >&2
  exit 1
fi
if ! docker compose version >/dev/null 2>&1; then
  printf '%s\n' "Docker Compose v2 is required. Install the Docker Compose plugin, then retry." >&2
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  printf '%s\n' "Docker daemon is unavailable. Start Docker, then retry." >&2
  exit 1
fi

printf '%s\n' "Validating Docker Compose configuration..."
"${compose[@]}" config --quiet

printf '%s\n' "Building and starting the library stack..."
if ! "${compose[@]}" up --detach --build --wait --wait-timeout 180; then
  printf '%s\n' "Setup failed. Container status:" >&2
  "${compose[@]}" ps >&2 || true
  printf '%s\n' "Inspect logs with: docker compose -f $project_dir/docker-compose.yaml logs --tail=200" >&2
  exit 1
fi

printf '%s\n' "Checking gateway upstreams..."
gateway_ready=false
for _ in {1..20}; do
  if "${compose[@]}" exec --no-TTY frontend node -e '
    const paths = [
      "/health/readiness",
      "/api/v1/docs/users/swagger.json",
      "/api/v1/docs/books/docs/swagger.json",
      "/api/v1/docs/transactions/swagger.json",
    ];
    Promise.all(paths.map(async (path) => {
      const response = await fetch(new URL(path, "http://gateway:8000"));
      if (!response.ok) throw new Error(path + ": " + response.status);
    })).catch(() => process.exit(1));
  ' >/dev/null 2>&1; then
    gateway_ready=true
    break
  fi
  sleep 2
done
if [[ "$gateway_ready" != true ]]; then
  printf '%s\n' "Gateway upstream checks failed. Container status:" >&2
  "${compose[@]}" ps >&2 || true
  printf '%s\n' "Inspect logs with: docker compose -f $project_dir/docker-compose.yaml logs --tail=200" >&2
  exit 1
fi

printf '%s\n' \
  "" \
  "Library stack is ready." \
  "Application: http://localhost:3000" \
  "API gateway: http://localhost:8000" \
  "RabbitMQ UI: http://localhost:15672 (library / library_password)" \
  "Member login: member@library.com / password" \
  "Admin login: admin@library.com / password"
