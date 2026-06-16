#!/usr/bin/env bash
# Health check script for homelab services.
set -uo pipefail

check() {
  local name="$1"
  local cmd="$2"
  if eval "$cmd" >/dev/null 2>&1; then
    echo "✓ $name"
  else
    echo "✗ $name"
  fi
}

echo "== homelab status =="

# GPU status
if command -v nvidia-smi >/dev/null 2>&1; then
  nvidia-smi --query-gpu=name,memory.used,memory.total --format=csv,noheader
else
  echo "✗ GPU (nvidia-smi not found)"
fi

# Docker services
check "docker: rag-platform" "docker ps --format '{{.Names}}' | grep -q '^rag-platform$'"
check "docker: palette-vein" "docker ps --format '{{.Names}}' | grep -E -q '(^palettevein|palette)'"

# Host services
check "Ollama (:11434)" "curl -sf http://localhost:11434/api/tags >/dev/null"
check "nginx (:80)" "curl -sf -o /dev/null http://localhost:80"
check "cloudflared" "systemctl is-active --quiet cloudflared"

# Endpoints
check "endpoint: palettevein.luy869.net" "curl -sf -o /dev/null --max-time 8 https://palettevein.luy869.net"
check "endpoint: api-internal.luy869.net" "curl -sf -o /dev/null --max-time 8 https://api-internal.luy869.net"

echo
