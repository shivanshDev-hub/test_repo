package classifier

// inference_test.go — unit tests for the Age ONNX backend and the
// image → 512-d embedding preprocessor.
//
// The real-ONNX tests are opt-in: set ONNX_RUNTIME_LIB to the path of a
// libonnxruntime.so (and optionally MODEL_PATH_AGE to the model). Without the
// library they t.Skip so `go test ./...` stays green in CI.

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"testing"
)

// tinyPNG returns a deterministic 32×32 gradient PNG.
func tinyPNG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 8),
				G: uint8(y * 8),
				B: uint8((x + y) * 4),
				A: 255,
			})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode test PNG: %v", err)
	}
	return buf.Bytes()
}

// TestExtractFaceEmbedding checks the deterministic stand-in embedding that
// feeds age_mlp.onnx.
func TestExtractFaceEmbedding(t *testing.T) {
	emb, err := extractFaceEmbedding(tinyPNG(t))
	if err != nil {
		t.Fatalf("extractFaceEmbedding: %v", err)
	}
	if len(emb) != embeddingDim {
		t.Fatalf("expected %d dims, got %d", embeddingDim, len(emb))
	}

	var norm float64
	for _, v := range emb {
		norm += float64(v) * float64(v)
	}
	norm = math.Sqrt(norm)
	if math.Abs(norm-1.0) > 1e-4 {
		t.Fatalf("embedding not L2-normalized (norm=%f)", norm)
	}

	// Determinism: same input → same vector.
	again, err := extractFaceEmbedding(tinyPNG(t))
	if err != nil {
		t.Fatalf("second extractFaceEmbedding: %v", err)
	}
	for i := range emb {
		if emb[i] != again[i] {
			t.Fatalf("embedding not deterministic at index %d", i)
		}
	}

	// No content → must not zero out the entire vector for realistic images;
	// and garbage input must produce an error.
	if _, err := extractFaceEmbedding([]byte("not an image")); err == nil {
		t.Fatal("expected error for undecodable bytes")
	}
	if _, err := extractFaceEmbedding(nil); err == nil {
		t.Fatal("expected error for empty bytes")
	}
}

// TestONNXBackendInfer loads the real age_mlp.onnx through onnxruntime_go and
// checks that Infer returns a probability in [0,1]. Skipped when the runtime
// library is not available (CI / environments without libonnxruntime).
func TestONNXBackendInfer(t *testing.T) {
	libPath := os.Getenv("ONNX_RUNTIME_LIB")
	if libPath == "" {
		t.Skip("ONNX_RUNTIME_LIB not set; skipping real ONNX inference test")
	}

	modelPath := os.Getenv("MODEL_PATH_AGE")
	if modelPath == "" {
		modelPath = "../weights/age_mlp.onnx"
	}
	if _, err := os.Stat(modelPath); err != nil {
		t.Skipf("model %s not present: %v", modelPath, err)
	}

	b, err := NewONNXBackend("age_test", modelPath)
	if err != nil {
		t.Fatalf("NewONNXBackend: %v", err)
	}
	defer b.close()

	score, err := b.Infer(tinyPNG(t))
	if err != nil {
		t.Fatalf("Infer: %v", err)
	}
	if score < 0 || score > 1 {
		t.Fatalf("p_minor out of range: %f", score)
	}
	t.Logf("age model p_minor for test image: %.6f", score)

	// Two identical calls must agree (session input handling is correct).
	score2, err := b.Infer(tinyPNG(t))
	if err != nil {
		t.Fatalf("second Infer: %v", err)
	}
	if score != score2 {
		t.Fatalf("inference not deterministic: %f != %f", score, score2)
	}
}

// TestONNXBackendEmptyInput guards the error path on empty bytes.
func TestONNXBackendEmptyInput(t *testing.T) {
	b, err := NewONNXBackend("age_test", "/models/age_mlp.onnx")
	if err != nil {
		t.Fatalf("NewONNXBackend: %v", err)
	}
	defer b.close()

	if _, err := b.Infer(nil); err == nil {
		t.Fatal("expected error for empty image bytes")
	}
}

// TestExtractImageTensor checks the CNN (Vector B) preprocessor: it must
// return exactly 3*size*size floats, be deterministic, and scale into the
// normalized range used by ImageNet-style models.
func TestExtractImageTensor(t *testing.T) {
	const size = 32
	tensor, err := extractImageTensor(tinyPNG(t), size, imagenetMean, imagenetStd)
	if err != nil {
		t.Fatalf("extractImageTensor: %v", err)
	}
	if len(tensor) != 3*size*size {
		t.Fatalf("expected %d floats, got %d", 3*size*size, len(tensor))
	}

	again, err := extractImageTensor(tinyPNG(t), size, imagenetMean, imagenetStd)
	if err != nil {
		t.Fatalf("second extractImageTensor: %v", err)
	}
	for i := range tensor {
		if tensor[i] != again[i] {
			t.Fatalf("image tensor not deterministic at index %d", i)
		}
	}

	if _, err := extractImageTensor([]byte("garbage"), size, imagenetMean, imagenetStd); err == nil {
		t.Fatal("expected error for undecodable bytes")
	}
}

// TestExtractPoseKeypoints checks the Vector C stand-in: exactly
// poseKeypointDim floats, each coordinate/confidence clamped to [0,1], and
// deterministic for identical input.
func TestExtractPoseKeypoints(t *testing.T) {
	kps, err := extractPoseKeypoints(tinyPNG(t))
	if err != nil {
		t.Fatalf("extractPoseKeypoints: %v", err)
	}
	if len(kps) != poseKeypointDim {
		t.Fatalf("expected %d keypoint features, got %d", poseKeypointDim, len(kps))
	}
	for i, v := range kps {
		if v < 0 || v > 1 {
			t.Fatalf("keypoint feature %d out of range: %f", i, v)
		}
	}

	again, err := extractPoseKeypoints(tinyPNG(t))
	if err != nil {
		t.Fatalf("second extractPoseKeypoints: %v", err)
	}
	for i := range kps {
		if kps[i] != again[i] {
			t.Fatalf("pose keypoints not deterministic at index %d", i)
		}
	}
}

// TestModelSpecs ensure every vector spec is well-formed (names, shapes,
// non-nil preprocessor) before any session is attempted.
func TestModelSpecs(t *testing.T) {
	specs := []ONNXModelSpec{AgeModelSpec()}

	nud, err := NudityModelSpec(224)
	if err != nil {
		t.Fatalf("NudityModelSpec: %v", err)
	}
	specs = append(specs, nud)

	beh, err := BehaviorModelSpec()
	if err != nil {
		t.Fatalf("BehaviorModelSpec: %v", err)
	}
	specs = append(specs, beh)

	for _, s := range specs {
		if err := s.validate(); err != nil {
			t.Fatalf("spec validate: %v", err)
		}
	}

	if _, err := NudityModelSpec(0); err == nil {
		t.Fatal("expected error for non-positive nudity model size")
	}
}

// TestNewONNXModelBackendBadSpec ensures the constructor rejects malformed
// specs without touching the filesystem or runtime.
func TestNewONNXModelBackendBadSpec(t *testing.T) {
	bad := ONNXModelSpec{
		InputName:  "x",
		OutputName: "y",
		InputShape: []int64{},
		Preprocess: func(_ []byte) ([]float32, error) { return nil, nil },
	}
	if _, err := NewONNXModelBackend("bad", "/tmp/nope.onnx", bad); err == nil {
		t.Fatal("expected error for empty input shape")
	}

	missingPP := ONNXModelSpec{
		InputName:  "x",
		OutputName: "y",
		InputShape: []int64{1, 1},
	}
	if _, err := NewONNXModelBackend("bad", "/tmp/nope.onnx", missingPP); err == nil {
		t.Fatal("expected error for missing preprocessor")
	}
}
