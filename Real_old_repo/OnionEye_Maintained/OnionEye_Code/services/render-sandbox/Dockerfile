# syntax=docker/dockerfile:1
#
# OnionEye — Render Sandbox
# Air-gapped Chromium renderer that produces real pixel screenshots from
# deeply sanitized HTML.  JS-disabled, network-isolated, non-root.
#
# ── Build stage ─────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY internal/ internal/
COPY cmd/ cmd/

# Hermetic: build from the committed vendor/ tree (no network downloads).
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod=vendor -trimpath -ldflags="-s -w" -o /bin/sandbox ./cmd/sandbox

# Small static healthcheck binary (stdlib only, no vendor needed).
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /bin/healthcheck ./cmd/healthcheck

# ── Runtime stage ──────────────────────────────────────────────────────────
FROM debian:bookworm-slim

# wget and ca-certificates are intentionally excluded: the container has no
# egress so TLS roots serve no purpose, and wget is a recon/deviator tool
# with no place in the runtime image.  The healthcheck binary replaces the
# former wget-based probe.
RUN apt-get update && apt-get install -y --no-install-recommends \
        chromium \
        fonts-liberation \
        fonts-noto-cjk \
        libnss3 \
        libatk1.0-0 \
        libatk-bridge2.0-0 \
        libcups2 \
        libdrm2 \
        libxkbcommon0 \
        libxcomposite1 \
        libxdamage1 \
        libxrandr2 \
        libgbm1 \
        libpango-1.0-0 \
        libcairo2 \
        libasound2 \
        libatspi2.0-0 \
    && rm -rf /var/lib/apt/lists/*

# The render service runs as the non-root `sandbox` user. -m creates its home
# directory; Chromium 152 derives the crashpad database path from $HOME, and
# without a writable home every render aborts with
# "chrome_crashpad_handler: --database is required" followed by a SIGTRAP core
# dump (manual verification: HOME=/tmp + --no-sandbox renders fine).
RUN groupadd -r sandbox && useradd -r -m -d /home/sandbox -g sandbox -G audio,video sandbox \
    && chown -R sandbox:sandbox /home/sandbox

ENV HOME=/home/sandbox

COPY --from=builder /bin/sandbox /bin/sandbox
COPY --from=builder /bin/healthcheck /bin/healthcheck

USER sandbox

LABEL org.opencontainers.image.title="onioneye-render-sandbox" \
      org.opencontainers.image.description="Air-gapped Chromium screenshot renderer (Team 2)" \
      org.opencontainers.image.source="https://github.com/onioneye/onioneye" \
      org.opencontainers.image.version="0.1.0"

EXPOSE 8084

HEALTHCHECK --interval=30s --timeout=5s --start-period=15s --retries=3 \
    CMD ["/bin/healthcheck"]

CMD ["/bin/sandbox"]