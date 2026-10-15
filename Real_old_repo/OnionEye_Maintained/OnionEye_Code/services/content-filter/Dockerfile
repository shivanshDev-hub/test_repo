# syntax=docker/dockerfile:1
#
# OnionEye — Content Filter (Team 2, Phase 8A)
# ONNX ML Abusive-Content Classifier. Runs age_mlp.onnx → p_minor through the
# Stage 1→4 pipeline behind a stateless Gin HTTP service.
# Zero-trust alignment: non-root, no capabilities, no-new-privileges.

# ── ONNX Runtime stage — official v1.29.0 shared library (matches the
# ── onnxruntime_go C API headers used at compile time = ORT_API_VERSION 29) ──
FROM debian:bookworm-slim AS onnxruntime
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && curl -fsSL -o /tmp/ort.tgz \
        https://github.com/microsoft/onnxruntime/releases/download/v1.29.0/onnxruntime-linux-x64-1.29.0.tgz \
    && tar -xzf /tmp/ort.tgz -C /tmp \
        onnxruntime-linux-x64-1.29.0/lib/libonnxruntime.so.1.29.0 \
    && mv /tmp/onnxruntime-linux-x64-1.29.0/lib/libonnxruntime.so.1.29.0 /tmp/libonnxruntime.so.1.29.0 \
    && ln -s /tmp/libonnxruntime.so.1.29.0 /tmp/libonnxruntime.so \
    && rm -rf /tmp/ort.tgz /tmp/onnxruntime-linux-x64-1.29.0 \
    && apt-get remove -y curl && apt-get autoremove -y && rm -rf /var/lib/apt/lists/*

# ── Build stage ────────────────────────────────────────────────────────────────
# The runtime image is Debian (glibc) and the official libonnxruntime.so is
# GLIBC-only, so the builder must ALSO be glibc: an Alpine (musl) binary cannot
# exec inside the Debian runtime ("no such file or directory" on the ELF
# interpreter) and cannot load the GLIBC-built ORT shared library either.
FROM golang:1.25-bookworm AS builder

RUN apt-get update \
    && apt-get install -y --no-install-recommends gcc libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

# Hermetic: build from the committed vendor/ tree (no network module downloads).
# onnxruntime_go is cgo-based, so CGO stays enabled here.
RUN CGO_ENABLED=1 go build -mod=vendor -trimpath \
    -ldflags="-s -w" -o /bin/content-filter ./cmd/filter

# ── Runtime stage ──────────────────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN addgroup --system --gid 10002 filter \
    && adduser --system --uid 10002 --ingroup filter filter \
    && apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

COPY --from=onnxruntime /tmp/libonnxruntime.so.1.29.0 /usr/lib/
COPY --from=onnxruntime /tmp/libonnxruntime.so /usr/lib/

COPY --from=builder /bin/content-filter /bin/content-filter

# Age model weights are baked in (read-only); mount over /models to override.
WORKDIR /models
COPY weights/ /models/

USER filter

ENV USE_MOCK_MODELS=true \
    MODEL_PATH_AGE=/models/age_mlp.onnx \
    ONNX_RUNTIME_LIB=/usr/lib/libonnxruntime.so \
    PORT=8082

EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -fsS http://127.0.0.1:8082/healthz || exit 1

CMD ["/bin/content-filter"]