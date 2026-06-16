#!/usr/bin/env bash
# Stop the containerised stack. Does NOT touch host services.
set -euo pipefail
cd "$(dirname "$0")/.."
docker compose down
