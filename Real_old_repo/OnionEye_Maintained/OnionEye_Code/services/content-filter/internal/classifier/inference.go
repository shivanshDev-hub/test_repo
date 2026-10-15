// Package classifier implements the Phase 8A ML Abusive-Content Classifier.
//
// Architecture (from Architecture_Subsection_3.md §1.3):
//
//	Vector A — Age      : InsightFace backbone → 512-d face embedding → p_minor
//	Vector B — Nudity   : NudeNet v2 / MobileNet-style CNN (ONNX)     → p_nudity
//	Vector C — Behavior : YOLOv8n-Pose + custom 2-layer MLP           → p_exploitative
//
// All three vectors run through the SAME onnxruntime session wrapper
// (ONNXBackend); only the ONNXModelSpec — tensor names, shapes, and the
// raw-image -> tensor preprocessor — differs per vector. During CI and the
// demo MockInferencer is used instead so the full pipeline is exercised
// without requiring the actual model files.
package classifier

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	ortgo "github.com/yalue/onnxruntime_go"
)

// Inferencer is the contract every model backend must satisfy.
// Keeping this as an interface lets the pipeline code stay identical whether
// the backing store is a real ONNX session or a deterministic stub.
type Inferencer interface {
	// Name returns a human-readable label used in logs and audit records.
	Name() string

	// Infer runs the model over the supplied image bytes and returns a
	// probability in [0.0, 1.0].  Implementations MUST be safe for
	// concurrent calls from multiple goroutines (Stage 2 calls the nudity
	// and behavior models in parallel).
	Infer(imageBytes []byte) (float64, error)
}

// ---------------------------------------------------------------------------
// Model specs — per-vector contracts
// ---------------------------------------------------------------------------

// PreprocessFunc turns raw image bytes into the flat float32 input tensor the
// model expects. The returned slice length must equal the product of
// InputShape. Implementations are pure and must be goroutine-safe.
type PreprocessFunc func(imageBytes []byte) ([]float32, error)

// ONNXModelSpec describes one model contract: input/output tensor names, the
// input tensor shape, an optional output shape (default {1,1}), and the
// raw-image → tensor preprocessor.
//
// Keep tensor names in sync with the exporter:
//
//	Vector A — input "face_embedding" [1,512],          output "p_minor"       [1,1]
//	Vector B — input "input"   [1,3,size,size],         output "sigmoid"       [1,1]
//	Vector C — input "keypoints"  [1,poseKeypointDim],  output "p_exploitative"  [1,1]
type ONNXModelSpec struct {
	InputName  string
	OutputName string

	// InputShape is the model's declared wrapped-shaped input (e.g. {1,512})
	// and doubles as the onnxruntime_go allocated-tensor shape.
	InputShape []int64

	// OutputShape is the expected output tensor shape; {1,1} is the common
	// single-probability shape shared by all three vectors today.
	OutputShape []int64

	// Preprocess converts image bytes into the flat input tensor. It must
	// return exactly prod(InputShape) floats and is called per-Infer.
	Preprocess PreprocessFunc
}

// validate sanity-checks the spec so downstream errors surface at startup.
func (s ONNXModelSpec) validate() error {
	if s.InputName == "" || s.OutputName == "" {
		return fmt.Errorf("model spec: input/output tensor names must not be empty")
	}
	if len(s.InputShape) == 0 {
		return fmt.Errorf("model spec %q: input shape must not be empty", s.InputName)
	}
	n := int64(1)
	for _, d := range s.InputShape {
		if d <= 0 {
			return fmt.Errorf("model spec %q: input shape dims must be positive, got %v", s.InputName, s.InputShape)
		}
		n *= d
	}
	if s.Preprocess == nil {
		return fmt.Errorf("model spec %q: preprocessor is nil", s.InputName)
	}
	return nil
}

// outputDims returns the output shape, defaulting to {1,1}.
func (s ONNXModelSpec) outputDims() []int64 {
	if len(s.OutputShape) == 0 {
		return []int64{1, 1}
	}
	return s.OutputShape
}

// AgeModelSpec is Vector A — the Age MLP over a 512-d face embedding.
func AgeModelSpec() ONNXModelSpec {
	return ONNXModelSpec{
		InputName:   "face_embedding",
		OutputName:  "p_minor",
		InputShape:  []int64{1, embeddingDim},
		OutputShape: []int64{1, 1},
		Preprocess:  extractFaceEmbedding,
	}
}

// NudityModelSpec is Vector B — a single-label CNN classifier over a resized,
// ImageNet-normalized NCHW tensor. modelSize is the model's expected spatial
// dimension (224 for MobileNet-style backbones).
func NudityModelSpec(modelSize int) (ONNXModelSpec, error) {
	if modelSize <= 0 {
		return ONNXModelSpec{}, fmt.Errorf("nudity model spatial size must be positive, got %d", modelSize)
	}
	return ONNXModelSpec{
		InputName:   "input",
		OutputName:  "sigmoid",
		InputShape:  []int64{1, 3, int64(modelSize), int64(modelSize)},
		OutputShape: []int64{1, 1},
		Preprocess: func(imageBytes []byte) ([]float32, error) {
			return extractImageTensor(imageBytes, modelSize, imagenetMean, imagenetStd)
		},
	}, nil
}

// BehaviorModelSpec is Vector C — a pose-MLP over normalized keypoints. The
// stand-in preprocessor emits 17 keypoints × (x, y, confidence) = 51 floats;
// this is the ONLY block that must change when the real YOLOv8n-Pose backbone
// is swapped in.
func BehaviorModelSpec() (ONNXModelSpec, error) {
	return ONNXModelSpec{
		InputName:   "keypoints",
		OutputName:  "p_exploitative",
		InputShape:  []int64{1, poseKeypointDim},
		OutputShape: []int64{1, 1},
		Preprocess:  extractPoseKeypoints,
	}, nil
}

// ---------------------------------------------------------------------------
// ONNX backend (production)
// ---------------------------------------------------------------------------

// ONNXBackend wraps a loaded ONNX Runtime session backed by onnxruntime_go
// and is shared by all three vectors. The network is loaded lazily on the
// first Infer call and reused afterwards; all calls are serialized through mu
// because the C session isn't guaranteed thread-safe for concurrent Run()
// from separate goroutines.
//
// The onnxruntime shared library is located via the ONNX_RUNTIME_LIB
// environment variable (default: "libonnxruntime.so" on the loader path).
type ONNXBackend struct {
	name      string
	modelPath string
	spec      ONNXModelSpec
	mu        sync.Mutex // serializes session access

	session *ortgo.AdvancedSession
	input   *ortgo.Tensor[float32]
	output  *ortgo.Tensor[float32]
}

// NewONNXBackend constructs an Age-vector backend (backwards-compatible
// convenience wrapper for callers/tests).
func NewONNXBackend(name, modelPath string) (*ONNXBackend, error) {
	return NewONNXModelBackend(name, modelPath, AgeModelSpec())
}

// NewONNXModelBackend constructs a backend for the provided spec. It does not
// touch the filesystem or the runtime until Load or Infer.
func NewONNXModelBackend(name, modelPath string, spec ONNXModelSpec) (*ONNXBackend, error) {
	if modelPath == "" {
		return nil, fmt.Errorf("onnx backend %q: modelPath must not be empty", name)
	}
	if err := spec.validate(); err != nil {
		return nil, fmt.Errorf("onnx backend %q: %w", name, err)
	}
	return &ONNXBackend{name: name, modelPath: modelPath, spec: spec}, nil
}

func (b *ONNXBackend) Name() string { return b.name }

// Load eagerly builds the ONNX session so callers can validate model
// availability at startup and decide whether to fall back to a mock. It is
// idempotent; subsequent calls return nil once the session exists.
func (b *ONNXBackend) Load() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.ensureLoaded()
}

// ensureLoaded initializes the onnxruntime environment once and builds the
// session for modelPath. Callers must already hold b.mu.
func (b *ONNXBackend) ensureLoaded() error {
	if b.session != nil {
		return nil
	}

	// Process-global ONNX Runtime environment. The shared library path must be
	// set before InitializeEnvironment; it defaults to the loader's path so the
	// standard libonnxruntime.so (or a docker-provided one) just works.
	if !ortgo.IsInitialized() {
		libPath := os.Getenv("ONNX_RUNTIME_LIB")
		if libPath == "" {
			libPath = "libonnxruntime.so"
		}
		ortgo.SetSharedLibraryPath(libPath)
		if err := ortgo.InitializeEnvironment(); err != nil {
			return fmt.Errorf("initialize onnxruntime (%s): %w", libPath, err)
		}
	}

	// Pre-allocate input/output tensors; they must outlive the session.
	// The input tensor is sized from the model spec so any of the three
	// vectors (or a future model) is handled by this same routine.
	var inputFloats int64 = 1
	for _, d := range b.spec.InputShape {
		inputFloats *= d
	}
	input, err := ortgo.NewTensor(ortgo.Shape(b.spec.InputShape),
		make([]float32, inputFloats))
	if err != nil {
		return fmt.Errorf("allocate input tensor %q: %w", b.spec.InputName, err)
	}
	output, err := ortgo.NewEmptyTensor[float32](ortgo.Shape(b.spec.outputDims()))
	if err != nil {
		_ = input.Destroy()
		return fmt.Errorf("allocate output tensor %q: %w", b.spec.OutputName, err)
	}
	opts, err := ortgo.NewSessionOptions()
	if err != nil {
		_ = input.Destroy()
		_ = output.Destroy()
		return fmt.Errorf("create session options: %w", err)
	}

	session, err := ortgo.NewAdvancedSession(b.modelPath,
		[]string{b.spec.InputName},
		[]string{b.spec.OutputName},
		[]ortgo.Value{input},
		[]ortgo.Value{output},
		opts)
	if err != nil {
		_ = input.Destroy()
		_ = output.Destroy()
		return fmt.Errorf("create onnx session %q: %w", b.modelPath, err)
	}

	b.session = session
	b.input = input
	b.output = output
	return nil
}

// close releases the session and tensors. It is idempotent and safe to leave
// uncalled for the lifetime of the process; provided for tests and graceful
// shutdown paths.
func (b *ONNXBackend) close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.session != nil {
		_ = b.session.Destroy()
		b.session = nil
	}
	if b.input != nil {
		_ = b.input.Destroy()
		b.input = nil
	}
	if b.output != nil {
		_ = b.output.Destroy()
		b.output = nil
	}
	return nil
}

// Infer runs the ONNX session and returns the model's probability in
// [0.0, 1.0]. Implementations are required (by the Inferencer contract) to be
// safe for concurrent calls; the internal mutex serializes access to the C
// session while preprocessing and tensor writes stay per-call.
func (b *ONNXBackend) Infer(imageBytes []byte) (float64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(imageBytes) == 0 {
		return 0.0, fmt.Errorf("onnx backend %q: empty image bytes", b.name)
	}

	if err := b.ensureLoaded(); err != nil {
		return 0.0, fmt.Errorf("onnx backend %q: %w", b.name, err)
	}

	feats, err := b.spec.Preprocess(imageBytes)
	if err != nil {
		return 0.0, fmt.Errorf("onnx backend %q: preprocess input %q: %w", b.name, b.spec.InputName, err)
	}
	if len(feats) != len(b.input.GetData()) {
		return 0.0, fmt.Errorf("onnx backend %q: preprocessor returned %d floats, model expects %d",
			b.name, len(feats), len(b.input.GetData()))
	}

	// Copy the features into the tensor's backing slice; the session reads the
	// same memory on Run().
	copy(b.input.GetData(), feats)
	if err := b.session.Run(); err != nil {
		return 0.0, fmt.Errorf("onnx backend %q: run session: %w", b.name, err)
	}

	data := b.output.GetData()
	if len(data) == 0 {
		return 0.0, fmt.Errorf("onnx backend %q: empty %s output", b.name, b.spec.OutputName)
	}
	p := float64(data[0])
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	return p, nil
}

// ---------------------------------------------------------------------------
// Mock backend (CI / demo / unit tests)
// ---------------------------------------------------------------------------

// MockInferencer returns a fixed score for every call.  It is goroutine-safe.
// The score is the literal value supplied at construction — no randomness —
// so tests can write deterministic assertions against the threshold logic.
type MockInferencer struct {
	name  string
	score float64
}

// NewMock creates a mock that always returns score regardless of image content.
// Panics if score is outside [0, 1] to catch misconfigured tests early.
func NewMock(name string, score float64) *MockInferencer {
	if score < 0 || score > 1 {
		panic(fmt.Sprintf("MockInferencer %q: score %.4f is outside [0,1]", name, score))
	}
	return &MockInferencer{name: name, score: score}
}

func (m *MockInferencer) Name() string { return m.name }

func (m *MockInferencer) Infer(_ []byte) (float64, error) {
	return m.score, nil
}

// ---------------------------------------------------------------------------
// Noisy mock (fuzzing / integration smoke tests)
// ---------------------------------------------------------------------------

// NoisyMock adds Gaussian jitter around a base score so integration tests can
// exercise the full [0,1] range without coupling to a specific seed.
type NoisyMock struct {
	name   string
	base   float64
	stddev float64
	mu     sync.Mutex
	rng    *rand.Rand
}

// NewNoisyMock creates a mock with base mean and Gaussian stddev.
func NewNoisyMock(name string, base, stddev float64, seed int64) *NoisyMock {
	return &NoisyMock{
		name:   name,
		base:   base,
		stddev: stddev,
		rng:    rand.New(rand.NewSource(seed)), //nolint:gosec — intentionally non-crypto
	}
}

func (n *NoisyMock) Name() string { return n.name }

func (n *NoisyMock) Infer(_ []byte) (float64, error) {
	n.mu.Lock()
	v := n.base + n.rng.NormFloat64()*n.stddev
	n.mu.Unlock()
	// clamp to [0, 1]
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	return v, nil
}
