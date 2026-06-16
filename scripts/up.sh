#!/usr/bin/env bash
# Start the containerised stack (chat_bot + palette-vein) via the root compose.
# Host services (Ollama, CLIP, nginx, cloudflared) are managed separately — see services/ and proxy/.
set -euo pipefail
cd "$(dirname "$0")/.."
docker compose up -d --build
echo "Stack started. Run ./scripts/status.sh to verify."
