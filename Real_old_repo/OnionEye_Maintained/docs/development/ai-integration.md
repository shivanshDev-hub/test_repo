# OnionEye: AI Integration & Setup Guide

This document outlines how to download, configure, and test the AI models and vector analytics engines across the OnionEye pipeline.

---

## 1. Ingestion API Models (Phases 6 & 7)
These models handle the input safety guarding and the seed term expansion. They are decoupled from the main Go process and served locally via the `ollama` container to manage GPU/CPU resources cleanly.

### A. Initial Seed Expander Model (Qwen2.5-3B)
- **Purpose:** Expands target search terms into typosquats, slang, and aliases using structured grammar constraints.
- **Container:** `ollama` (Exposed natively on internal port `11434`)
- **Environment Variable Setup:** Defined in `docker-compose.yml` under `ingestion-api` as `EXPANDER_MODEL=qwen2.5:3b`.
- **How to Add/Initialize:**
  1. Ensure the Ollama container is running with the AI profile: 
     ```bash
     docker compose --profile ai up -d ollama
     ```
  2. Pull the model weights locally inside the container: 
     ```bash
     docker compose exec ollama ollama pull qwen2.5:3b
     ```

### B. Input/Output Guard Model (LLaMA-Guard3-1B)
- **Purpose:** Acts as a dual-filter to screen malicious analyst intent (Phase 6b) and sanitize the AI-generated expansions from Qwen (Phase 8).
- **Environment Variable Setup:** Defined in `docker-compose.yml` as `GUARD_MODEL=llama-guard3:1b`.
- **How to Add/Initialize:**
  ```bash
  docker compose exec ollama ollama pull llama-guard3:1b
  ```

### C. Fail-Open & Offline Fallback Mechanism
- If the Ollama runtime is not running or unreachable, `ingestion-api` automatically falls back to:
  1. **Tier-1 Heuristic / Regex Blocklist:** Sub-5ms fast screening for dangerous SQL, command injection, and system payloads.
  2. **Pass-Through Default:** Marks the guard verdict as `PASS_NO_LLM` and processes raw unexpanded seed terms without throwing HTTP 500 errors.

### Testing the Ingestion Models
To verify the models are loaded and outputting correctly, you can query the Ollama API directly from your terminal:
```bash
# Test the Expander Output
curl -X POST http://localhost:11434/api/generate -d '{
  "model": "qwen2.5:3b",
  "prompt": "Expand the dark web vendor handle: shadow_vendor",
  "stream": false
}'

# Test the Guard Output
curl -X POST http://localhost:11434/api/generate -d '{
  "model": "llama-guard3:1b",
  "prompt": "Evaluate safety: PASS or BLOCK. Input: DROP TABLES",
  "stream": false
}'
```

---

## 1.5 Dashboard AI Controls & AI Log Feeds

The Tactical Web Console exposes runtime control of the AI engine directly from the header bar:

### AI Mode Indicator & Restart Buttons
- **`AI: ON/OFF` badge** — Reflects whether the Ollama engine container is currently up (polled from `GET /api/v1/dashboard/ai/status` every 10s).
- **`RESTART WITH AI` button** — Confirms a warning modal, then tells the backend to run `docker compose --profile ai up -d --force-recreate` (via `POST /api/v1/dashboard/ai/restart {enabled:true}`).
- **`RESTART NO AI` button** — Same flow with `{enabled:false}`: plain `up -d`, then the `ollama` container is explicitly stopped so no AI processing happens.

Both buttons warn that the whole stack restarts in **~15 seconds** and that the page must be **reloaded afterward** to reconnect to the freshly started services.

> **Backend plumbing:** `dashboard-api` mounts the Docker engine socket (read-only) *and* the `docker-compose.yml` file, and its image now ships `docker-cli` + `docker-cli-compose-plugin`. The `/ai/restart` handler shells out to `docker compose` from inside the container using the same project name as the running stack.

### AI Log Feeds (SYS LOGS tab)
Three virtual feeds surface AI activity in the SYS LOGS → component-chips row:

| Feed | Meaning |
|------|---------|
| **`ai`** | Combined view of every AI-relevant line from `ingestion-api` + `ollama`. |
| **`ai-seed`** | Initial **seed / term-expansion** stage (Qwen2.5-3B) — the first AI touchpoint of every investigation. |
| **`ai-guard`** | **CSAM / safety-guard block** stage (LLaMA-Guard3-1B) — the dual-filter that screens analyst input and sanitizes AI output. |

**These feeds are disabled outside AI mode:** when the Ollama engine is not running, the feeds (and their chips) are hidden from the SYS LOGS panel entirely. Only when the tool runs with the `--profile ai` stack (Ollama up) do the `ai`, `ai-seed`, and `ai-guard` feeds become selectable.

> **Backend plumbing:** `GET /api/v1/dashboard/logs?service=ai-seed|ai-guard|ai` returns lines filtered from the real `ingestion-api` and `ollama` container logs using stage-specific keywords. The combined "all logs" response only includes the AI feeds while Ollama is running.

---

## 2. Stylometric 384-D Vector Analytics (Phase 9)
The correlation engine analyzes forum posts, listing descriptions, and communications using stylometric dense vectors (384 dimensions) to link disparate threat actor personas across disconnected hidden services.

### A. Vector Architecture & Gating
- **Dense Embedding:** 384-dimensional dense vectors representing sentence structure, punctuation frequencies, vocabulary richness, and syntax patterns.
- **Storage & Indexing:** Stored in PostgreSQL using `pgvector` with HNSW (Hierarchical Navigable Small World) cosine distance indexing:
  ```sql
  CREATE INDEX IF NOT EXISTS idx_scans_stylometry 
  ON scans USING hnsw (stylometry_vector vector_cosine_ops);
  ```
- **Statistical Sample Size Gating:** Authorship attribution is statistically unreliable on short snippets. OnionEye enforces a **150-word minimum threshold** (`MinStylometryWords = 150`). Comparisons with fewer than 150 words receive an automatic $0.25\times$ penalty to prevent false-positive links.

### B. Testing Stylometry & Vector Scoring
You can test the vector similarity endpoint directly on `correlation-engine`:
```bash
curl -X POST http://localhost:8083/evaluate/stylometry \
  -H "Content-Type: application/json" \
  -d '{
    "vector_a": [0.02, 0.87, 0.15, 0.42],
    "vector_b": [0.03, 0.85, 0.14, 0.40],
    "sample_words_a": 320,
    "sample_words_b": 280
  }'
```

**Expected Response:**
```json
{
  "cosine_similarity": 0.9984,
  "confidence": 0.9984,
  "gated": false,
  "note": "sample size verified (>150 words)"
}
```

---

## 3. CSAM Filtration AI Models (Phase 8A)
These 3 models form the secondary deep ML Abusive-Content Classifier, acting as defense-in-depth for visual media that has bypassed the Phase 5 hash-based shield.

### A. Age Vector Model (Ordinal Regression)
- **Purpose:** Detects minor status. If probability $p_{\text{minor}} < 0.40$, the pipeline short-circuits and skips further extraction.
- **Model Stack:** ArcFace / InsightFace MS1MV2 backbone + custom 2-layer MLP head.
- **Placement:** Place ONNX weights under: 
  `OnionEye_Code/services/content-filter/weights/age_mlp.onnx`

### B. Nudity Vector Model (Knowledge Distilled)
- **Purpose:** Evaluates explicit visual elements (safe vs explicit).
- **Model Stack:** Compact MobileNetV3-Small / EfficientNet-Lite exported to ONNX (FP16/INT8).
- **Placement:** Place ONNX weights under: 
  `OnionEye_Code/services/content-filter/weights/nudity_mobilenet.onnx`

### C. Behavior & Pose Vector Model
- **Purpose:** Evaluates suggestive/exploitative poses via skeletal landmarks, distinct from harmless imagery.
- **Model Stack:** Frozen YOLOv8n-pose (to extract the 17 keypoints) fed into a lightweight MLP.
- **Placement:**
  - Base YOLO weights: `OnionEye_Code/services/content-filter/weights/yolov8n-pose.pt`
  - MLP weights: `OnionEye_Code/services/content-filter/weights/pose_mlp.onnx`

### Testing the ML Content Filter Endpoint
The `content-filter` service (`docker compose up -d content-filter`, or
`go run ./cmd/filter` inside `services/content-filter`) exposes:

- `POST /classify` — single image (multipart field `image`, optional `src`/`actor_id`)
- `POST /classify/job` — multiple images (`raw_b64`, standard base64)
- `GET /healthz` — liveness

```bash
curl -X POST http://localhost:8082/classify \
  -F "image=@/path/to/face.jpg" \
  -F "src=http://mock-market-1.onion/listing/8832/image.jpg"
# → { "stage_scores": { "p_minor": 0.48148, ... }, "classifier_score": 0.23217, ... }
```

Runtime knobs (see `services/content-filter/weights/README.md`):

| env | default | meaning |
|---|---|---|
| `USE_MOCK_MODELS` | unset → mocks | `"false"` attempts real ONNX loads per vector (otherwise pure mocks) |
| `MODEL_PATH_AGE` | `/models/age_mlp.onnx` | Vector A — age MLP path |
| `MODEL_PATH_NUDITY` | `/models/nudity_mobilenet.onnx` | Vector B — nudity CNN path |
| `MODEL_PATH_BEHAVIOR` | `/models/pose_mlp.onnx` | Vector C — pose MLP path |
| `NUDITY_MODEL_SIZE` | `224` | Vector B input spatial size (e.g. `224` → `[1,3,224,224]`); must match the exported model |
| `ONNX_RUNTIME_LIB` | `libonnxruntime.so` | onnxruntime shared library |

### Per-Vector ONNX Specs

All three vectors run through the **same** `ONNXBackend` wrapper; only the
`ONNXModelSpec` differs — input/output tensor names, tensor shape, and the
raw-image → tensor preprocessor:

| Vector | Input tensor | Shape | Output tensor | Preprocessor |
|---|---|---|---|---|
| A — Age | `face_embedding` | `[1,512]` | `p_minor` | `extractFaceEmbedding` (deterministic 512-d stand-in) |
| B — Nudity | `input` | `[1,3,224,224]` NCHW | `sigmoid` | `extractImageTensor` (bilinear resize + ImageNet norm) |
| C — Behavior | `keypoints` | `[1,51]` | `p_exploitative` | `extractPoseKeypoints` (17-band → pseudo-COCO stand-in) |

Keep these tensor names in sync with the exporter; a mismatch fails load
eagerly at startup and that vector falls back to its mock.

### Fail-Soft Mock Fallback

Loading is attempted per vector at startup. A vector uses the **real ONNX
model** only if all of these hold: `USE_MOCK_MODELS=false`, the file at the
vector's `MODEL_PATH_*` exists, and the onnxruntime session loads the model
with shapes matching its spec. Otherwise that single vector falls back to a
deterministic `MockInferencer` (fixed score per vector) so the service stays up
and the rest of the pipeline keeps running. **`USE_MOCK_MODELS` defaults to
`true`** (unset or anything-but-`false` → mock detection, which is what flags
content such as `p_minor=0.99 → held_pending_review`); set it to `false` on the
command line or via the compose/profile to attempt the real ONNX weights.
Confirm which mode each vector is in via the startup log lines (e.g.
`age: loaded real model` vs `age: model path not configured; using mock`).

> **Status (2026-09-09):** the ONNX backend is fully generalized to all three
> vectors with exported spec constructors (`AgeModelSpec`,
> `NudityModelSpec(size)`, `BehaviorModelSpec`) and per-vector preprocessors.
> `USE_MOCK_MODELS` now defaults to `true` (mock detection → content is held,
> i.e. `p_minor=0.99`), which is the out-of-the-box behavior. The bundled
> `age_mlp.onnx` ships with *synthesized placeholder weights*, so real-model
> `p_minor` values are demo-scale until the real trained export is dropped in.
> Nudity and Behavior vectors use deterministic mocks pending their ONNX
> artifacts. To activate the real ONNX path: drop `nudity_mobilenet.onnx` /
> `pose_mlp.onnx` into `services/content-filter/weights/`, set
> `USE_MOCK_MODELS=false`, and verify the startup log confirms real sessions
> for all three vectors.
