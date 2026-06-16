# Ollama (host service)

Ollama runs on the **host machine** (not containerised) for GPU access on port **11434**, bound to 0.0.0.0.

## Connectivity

- **From containers**: reachable at `http://172.17.0.1:11434`
- **From the host**: reachable at `http://localhost:11434`
- The chat_bot service connects via the `OLLAMA_HOST` environment variable (set to `http://172.17.0.1:11434`)

## Models

The following models are available:
- **qwen3.6** — general-purpose LLM
- **gemma4** — lightweight LLM
- **bge-m3** — embeddings model
- Cloud models — via API fallback

## Configuration

GPU assignment is controlled via environment variables in `apps/chat_bot/.env`:
- `GPU_DEVICES=0,1` — which GPUs Ollama can use
- `OLLAMA_PARALLEL=2` — parallel request handling

## How to Run

Start Ollama on the host:
```bash
ollama serve
```

Typically managed as a systemd service for auto-startup:
```bash
systemctl start ollama
systemctl status ollama
```

## Why on the Host?

Ollama runs on the host machine for two reasons:
1. **GPU access** — direct hardware access for inference
2. **Dual-use** — the host machine itself may need Ollama for desktop use
