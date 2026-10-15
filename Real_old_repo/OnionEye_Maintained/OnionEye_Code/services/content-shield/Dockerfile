# syntax=docker/dockerfile:1
#
# OnionEye — Content Safety Shield (Team 2)
# Perceptual-hash CSAM/bystander filter. Runs on raw bytes before any decoding.
# Zero-trust: non-root, read-only root.

# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY . .

# Hermetic: build from the committed vendor/ tree (no network module downloads).
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod=vendor -trimpath -ldflags="-s -w" -o /bin/shield ./cmd/shield

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN addgroup -S -g 10002 shield && adduser -S -u 10002 -G shield shield

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY --from=builder /bin/shield /bin/shield

USER shield

LABEL org.opencontainers.image.title="onioneye-content-shield" \
      org.opencontainers.image.description="CSAM/bystander perceptual-hash filter (Team 2)" \
      org.opencontainers.image.source="https://github.com/onioneye/onioneye" \
      org.opencontainers.image.version="0.1.0"

EXPOSE 8082

# Hash set is mounted read-only at /data/hashes (see compose)
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8082/health || exit 1

CMD ["/bin/shield"]
