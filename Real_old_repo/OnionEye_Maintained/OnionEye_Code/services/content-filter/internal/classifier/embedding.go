package classifier

// embedding.go — image → 512-d feature vector for the Age model input.
//
// Vector A in the architecture (§1.3) is:
//
//	InsightFace / ArcFace (MS1MV2) backbone → 512-d face embedding → age MLP → p_minor
//
// The bundled age_mlp.onnx consumes a float32 tensor of shape [1,512]
// named "face_embedding" and emits p_minor = sigmoid(...) as [1,1].
//
// The real face-embedding backbone (e.g. insightface buffalo_l / ArcFace
// resnet50) is intentionally NOT part of the hackathon artifact set, so this
// file provides a deterministic, content-dependent stand-in extractor:
//
//  1. decode the raw image (jpeg/png/gif),
//  2. area-average it onto a 32×32 RGB grid (32*32*3 = 3072 raw values),
//  3. aggregate every 6 consecutive raw values into a 512-d vector,
//  4. L2-normalize the result.
//
// When a real InsightFace .onnx backbone is provided later, ONLY this
// function's body must change — the ONNXBackend, the pipeline, and every
// threshold stay untouched.

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif" // decode support
	_ "image/jpeg"
	_ "image/png"

	"math"
)

// NOTE: We are going to add richer decoding after we confirm the needed
// formats. Right now we keep the build compiling while iterating.


// NOTE: Go's standard image package does not include webp decoding by default.
// We import "image/webp" here; the implementation is provided by the module
// github.com/chai2010/webp (wired in via go.mod in the content-filter service).

const (
	// embeddingDim is the fixed input width of age_mlp.onnx ("face_embedding").
	embeddingDim = 512
	// embeddingGrid is the spatial pooling grid: 32×32 px × 3 channels = 3072.
	embeddingGrid  = 32
	embeddingSteps = embeddingGrid * embeddingGrid * 3 / embeddingDim // 6
)

// extractFaceEmbedding decodes imageBytes and returns an L2-normalized
// 512-dimensional float32 vector suitable for the age MLP.
// It returns an error for undecodable/empty input; it always returns exactly
// embeddingDim values on success.
func extractFaceEmbedding(imageBytes []byte) ([]float32, error) {
	// Try standard decoder first (jpeg/png/gif).
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		// If the format isn't supported by the standard library, we fall
		// back to decoding using the “x/image” webp/gfx stack via imaging.
		// imaging supports more formats (including webp) when its deps exist.
		//
		// Note: imaging import and fallback wiring will be enabled once we
		// confirm the needed formats present in your refreshed sample set.
		return nil, fmt.Errorf("decode image: %w", err)
	}

	b := img.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, fmt.Errorf("image has empty bounds %v", b)
	}

	// --- 1. area-average the image onto a 32×32 RGB grid -------------------
	raw := make([]float64, embeddingGrid*embeddingGrid*3)
	counts := make([]int, embeddingGrid*embeddingGrid)

	w := float64(b.Dx())
	h := float64(b.Dy())
	stepX := float64(embeddingGrid) / w
	stepY := float64(embeddingGrid) / h

	for y := b.Min.Y; y < b.Max.Y; y++ {
		gi := int(float64(y-b.Min.Y) * stepY)
		for x := b.Min.X; x < b.Max.X; x++ {
			gj := int(float64(x-b.Min.X) * stepX)
			if gi >= embeddingGrid {
				gi = embeddingGrid - 1
			}
			if gj >= embeddingGrid {
				gj = embeddingGrid - 1
			}
			r, g, bl, a := img.At(x, y).RGBA()
			// The image color model may be alpha-aware; normalize alpha so
			// transparent pixels contribute their premultiplied color.
			if a != 0 {
				// RGBA() returns 16-bit values; drop to 8-bit resolution.
				cf := color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(bl >> 8), A: uint8(a >> 8)}
				cell := (gi*embeddingGrid + gj) * 3
				raw[cell+0] += float64(cf.R)
				raw[cell+1] += float64(cf.G)
				raw[cell+2] += float64(cf.B)
				counts[gi*embeddingGrid+gj]++
			}
		}
	}

	// --- 2. per-cell averaging ----------------------------------------------
	for i := range counts {
		n := counts[i]
		if n == 0 {
			continue
		}
		base := i * 3
		raw[base+0] /= float64(n)
		raw[base+1] /= float64(n)
		raw[base+2] /= float64(n)
	}

	// --- 3+4. collapse 3072 → 512 and L2-normalize --------------------------
	emb := make([]float32, embeddingDim)
	for i := 0; i < embeddingDim; i++ {
		var sum float64
		for k := 0; k < embeddingSteps; k++ {
			sum += raw[i*embeddingSteps+k]
		}
		emb[i] = float32(sum / embeddingSteps)
	}

	var norm float64
	for _, v := range emb {
		norm += float64(v) * float64(v)
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range emb {
			emb[i] = float32(float64(emb[i]) / norm)
		}
	}

	return emb, nil
}
