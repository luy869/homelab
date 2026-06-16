# CLIP gRPC service (host service)

CLIP runs on the **host machine** (not containerised) for GPU access on port **50051** via gRPC.

## Connectivity

- **From palette-vein backend**: reachable at `host.docker.internal:50051`
- The Docker Compose file maps `host.docker.internal` to the host gateway via `extra_hosts: ["host.docker.internal:host-gateway"]`
- This allows containers to reach host services without hardcoding IP addresses

## Purpose

Provides image embedding generation for the palette-vein application, enabling semantic image search and similarity matching.

## Starting the Service

Start CLIP from the palette-vein repository:
```bash
cd <palette-vein-repo>/clip_service
./run-clip.sh
```

This script handles model loading and starts the gRPC server on port 50051.

## Why on the Host?

CLIP runs on the host machine for **direct GPU access**, ensuring fast embedding generation without container overhead.
