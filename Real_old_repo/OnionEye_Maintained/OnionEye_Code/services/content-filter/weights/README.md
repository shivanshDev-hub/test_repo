# content-filter model weights

Phase 8A ML Abusive-Content Classifier — one ONNX file per detection vector.
All three vectors share one `ONNXBackend` session wrapper; only the
`ONNXModelSpec` differs. Each file is **optional**: if it is absent, that
vector falls back to a deterministic mock (fail-soft).

| File | Vector | Input | Output |
|---|---|---|---|
| `age_mlp.onnx` | A — Age | `face_embedding` `[1,512]` → `p_minor` `[1,1]` | Age MLP over 512-d face embedding |
| `nudity_mobilenet.onnx` | B — Nudity | `input` `[1,3,{size},{size}]` NCHW → `sigmoid` `[1,1]` | MobileNet-class CNN classifier |
| `pose_mlp.onnx` | C — Behavior | `keypoints` `[1,51]` → `p_exploitative` `[1,1]` | MLP over 17-keypoint pose features |

Architecture reference: `docs/development/ai-integration.md` §3.

## Vector A — Age detection MLP (`age_mlp.onnx`)

A PyTorch-exported 2-layer MLP (Gemm → ReLU → Gemm → Sigmoid) that produces
`p_minor` from a 512-dimensional face embedding:

| | name | type | shape |
|---|---|---|---|
| input | `face_embedding` | float32 | `[1, 512]` |
| output | `p_minor` | float32 | `[1, 1]` |

## ⚠️ Placeholder weights

The original user artifact was exported with **external tensor data**
(`age_mlp.onnx.data`), and that data file was missing from the project folder.
Because ONNX models cannot run without their weight tensors, this repository
ships a **re-exported single-file model with deterministically generated
placeholder weights** (numpy `default_rng(20260909)`, scale 0.08).

Consequences:

- The graph topology, tensor names, and I/O contract are exactly the user's
  model; the pipeline integration is 100% real.
- Score semantics are **synthetic** until the real weights are dropped in.
  Treat any `p_minor` value as demo-scale, not trained output.

### Replacing with the real weights

When you obtain the original training export, restore the artifact pair so the
runtime loads YOUR numbers:

```bash
# Option A — keep external-data format:
cp age_mlp.onnx      services/content-filter/weights/age_mlp.onnx
cp age_mlp.onnx.data services/content-filter/weights/age_mlp.onnx.data
```

or re-export from PyTorch as a single self-contained file:

```python
import torch, onnx
m = torch.load("age_mlp.pt")          # your trained AgeMLP
dummy = torch.randn(1, 512)           # embedding input
torch.onnx.export(m, dummy, "age_mlp.onnx", input_names=["face_embedding"],
                  output_names=["p_minor"], opset_version=18)
```

The runtime only needs `MODEL_PATH_AGE` to point at a valid `age_mlp.onnx`
(single file or with its sibling `.data`), so either layout works.

## Vectors B & C — pending artifacts

- **`nudity_mobilenet.onnx`** — compact MobileNetV3-Small / EfficientNet-Lite
  exported to ONNX (FP16/INT8). Input name `input`, **NCHW**
  `[1,3,size,size]`. Keep `NUDITY_MODEL_SIZE` in sync with the export
  (default 224).
- **`pose_mlp.onnx`** — lightweight MLP over YOLOv8n-pose 17-keypoint features
  (flat `[1,51]` = 17×(x→y, confidence)×3 normalised values). Input name
  `keypoints`, output `p_exploitative`.

Until these files are dropped in, Vectors B and C log
`model ... not present ...; using mock` at startup and return the configured
deterministic mock scores.

## Runtime configuration

| env | defaults | purpose |
|---|---|---|
| `USE_MOCK_MODELS` | anything-but-`false` | `"false"` attempts the real ONNX loads (per-vector fail-soft, else mock) |
| `MODEL_PATH_AGE` | `/models/age_mlp.onnx` | Vector A path |
| `MODEL_PATH_NUDITY` | `/models/nudity_mobilenet.onnx` | Vector B path |
| `MODEL_PATH_BEHAVIOR` | `/models/pose_mlp.onnx` | Vector C path |
| `NUDITY_MODEL_SIZE` | `224` | Vector B spatial size (must match the export) |
| `ONNX_RUNTIME_LIB` | `libonnxruntime.so` | onnxruntime shared library |

Local run example:

```bash
export ONNX_RUNTIME_LIB=/tmp/ort/onnxruntime-linux-x64-1.29.0/lib/libonnxruntime.so
export MODEL_PATH_AGE=$PWD/services/content-filter/weights/age_mlp.onnx
export MODEL_PATH_NUDITY=$PWD/services/content-filter/weights/nudity_mobilenet.onnx
export MODEL_PATH_BEHAVIOR=$PWD/services/content-filter/weights/pose_mlp.onnx
export USE_MOCK_MODELS=false
go run ./cmd/filter
```

Missing files fall back to mocks; the startup log shows exactly which mode each
vector is in (e.g. `OK: Vector A/age loaded real model from ...` vs
`WARNING: Vector B/nudity: model /models/nudity_mobilenet.onnx not present ...;
using mock`).