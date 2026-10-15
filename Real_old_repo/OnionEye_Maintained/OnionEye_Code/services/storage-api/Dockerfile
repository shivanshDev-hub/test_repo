# syntax=docker/dockerfile:1
#
# OnionEye — Storage API (Team 2)
# The ONLY container holding DB credentials. All writes flow through this HTTP API.
# Zero-trust boundary: crawler/parser/shield containers never have DATABASE_URL.

# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /app

COPY . .

# Hermetic build from the committed vendor/ tree — no network module downloads.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -mod=vendor -trimpath -ldflags="-s -w" -o /bin/server ./cmd/server

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN addgroup -S -g 10003 store && adduser -S -u 10003 -G store store

# psql + certs for postgres connectivity
RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY --from=builder /bin/server /bin/server

USER store

LABEL org.opencontainers.image.title="onioneye-storage-api" \
      org.opencontainers.image.description="Postgres storage service API (Team 2, sole DB credential holder)" \
      org.opencontainers.image.source="https://github.com/onioneye/onioneye" \
      org.opencontainers.image.version="0.1.0"

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://127.0.0.1:8081/health || exit 1

CMD ["/bin/server"]
