# ── Build stage ────────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
# Hermetic build from the committed vendor/ tree — no network module downloads.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /bin/engine ./cmd/engine

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /bin/engine /bin/engine
EXPOSE 8083
CMD ["/bin/engine"]
