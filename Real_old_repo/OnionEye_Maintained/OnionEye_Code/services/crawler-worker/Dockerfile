# syntax=docker/dockerfile:1
#
# OnionEye — Crawler Worker Pool (Team 2)
# Zero-trust container: non-root, dropped capabilities, read-only root.
# Runtime controls (read_only, cap_drop, tmpfs, seccomp, init, limits) are
# enforced in docker-compose.yml.
#
# Runtime stage is `scratch` — the absolute least-privilege image:
#   * contains ONLY the static binary + CA bundle (system certs so https
#     egress — including .onion TLS — keeps verifying correctly);
#   * no shell, no package manager, no setuid utilities, no debugger ⇒ a
#     compromised worker has nothing to execute/install/persist with;
#   * fully static (CGO_ENABLED=0) so no glibc/zoneinfo/node_modules;
#   * Docker still bind-mounts /etc/resolv.conf + /etc/hosts, so the
#     pure-Go DNS resolver works unchanged.

# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY . .

# CA bundle is present in the golang:alpine base; make it deterministic.
RUN apk add --no-cache ca-certificates 2>/dev/null || true

# Static binary — no CGO, no runtime dynamic deps
# Hermetic: build from the committed vendor/ tree (no network module downloads).
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod=vendor -trimpath -ldflags="-s -w" -o /bin/worker ./cmd/worker

# ── Runtime stage — scratch (minimal image, nothing to escalate with) ─────────
FROM scratch

# System CA bundle at the path Go's x509 SystemCertPool checks first —
# same concatenated PEM that alpine's ca-certificates-bundle ships.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /bin/worker /bin/worker

# Non-root, numeric UID (scratch has no /etc/passwd and none is needed here).
# 65532 == distroless nonroot convention; 65534 == nobody.
USER 65532:65532

# Metadata
LABEL org.opencontainers.image.title="onioneye-crawler-worker" \
      org.opencontainers.image.description="Tor crawler worker pool (Team 2) — scratch, non-root, all caps dropped" \
      org.opencontainers.image.source="https://github.com/onioneye/onioneye" \
      org.opencontainers.image.version="0.1.0"

# Worker is a daemon that reads jobs from Redis — no port to expose.
# It must NEVER hold DB credentials (zero-trust: writes go via storage-api).

CMD ["/bin/worker"]
