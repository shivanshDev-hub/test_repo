package classifier

// preprocess.go — raw image → model input tensor extractors used by the
// ONNX backend specs for Vectors B (Nudity CNN) and C (Behavior pose MLP).
//
// Vector B expects the standard CNN input layout: an ImageNet-normalized,
// NCHW-ordered float32 tensor, e.g. [1,3,224,224].
//
// Vector C expects a flat keypoint vector, [1, 51] = 17 COCO keypoints ×
// (x, y, confidence). The real backbone is YOLOv8n-Pose; until that ONNX is
// provided, extractPoseKeypoints returns a deterministic stand-in derived
// from the raw image so the pipeline is exercised end-to-end.

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif" // decode support
	_ "image/jpeg"
	_ "image/png"
)

// poseKeypointDim is the Vector C input width: 17 COCO-skeleton landmarks ×
// 3 values (normalized x, normalized y, confidence).
const poseKeypointDim = 17 * 3

// imagenetMean / imagenetStd are the standard normalization used by
// MobileNet / EfficientNet-style backbones exported from PyTorch/TFLite.
var (
	imagenetMean = [3]float64{0.485, 0.456, 0.406}
	imagenetStd  = [3]float64{0.229, 0.224, 0.225}
)

// extractImageTensor decodes imageBytes, bilinearly resizes it to size×size
// and returns an NCHW-ordered float32 vector of len 3*size*size normalized by
// mean/std. This is the input contract of virtually every MobileNet /
// EfficientNet CNN export (Vector B).
func extractImageTensor(imageBytes []byte, size int, mean, std [3]float64) ([]float32, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	b := img.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, fmt.Errorf("image has empty bounds %v", b)
	}

	// Bilinear sample the image onto a size×size grid. Source coordinates use
	// 0.5-offset sampling so destination pixel centers map without edge bias,
	// matching standard PyTorch/TFLite resize behavior.
	dst := make([]float64, size*size*3)
	for y := 0; y < size; y++ {
		sy := (float64(y)+0.5)/float64(size)*float64(b.Dy()) - 0.5
		y0 := int(sy)
		if y0 < 0 {
			y0 = 0
		}
		if y0 > b.Max.Y-1 {
			y0 = b.Max.Y - 1
		}
		cy := sy - float64(y0)
		if cy < 0 {
			cy = 0
		}
		for x := 0; x < size; x++ {
			sx := (float64(x)+0.5)/float64(size)*float64(b.Dx()) - 0.5
			x0 := int(sx)
			if x0 < 0 {
				x0 = 0
			}
			if x0 > b.Max.X-1 {
				x0 = b.Max.X - 1
			}
			cx := sx - float64(x0)
			if cx < 0 {
				cx = 0
			}

			x1 := x0 + 1
			if x1 > b.Max.X-1 {
				x1 = b.Max.X - 1
			}
			y1 := y0 + 1
			if y1 > b.Max.Y-1 {
				y1 = b.Max.Y - 1
			}

			r00, g00, bl00, _ := img.At(x0, y0).RGBA()
			r10, g10, bl10, _ := img.At(x0, y1).RGBA()
			r01, g01, bl01, _ := img.At(x1, y0).RGBA()
			r11, g11, bl11, _ := img.At(x1, y1).RGBA()

			w00 := (1 - cx) * (1 - cy)
			w10 := (1 - cx) * cy
			w01 := cx * (1 - cy)
			w11 := cx * cy

			// RGBA() returns 16-bit premultiplied values; blend then scale to
			// [0,1]. Opaque photos dominate here; ignoring alpha is the same
			// behavior the export-time preprocessing uses.
			to01 := func(v float64) float64 { return v / 0xFFFF }
			base := (y*size + x) * 3
			dst[base+0] = to01(w00*float64(r00) + w10*float64(r10) + w01*float64(r01) + w11*float64(r11))
			dst[base+1] = to01(w00*float64(g00) + w10*float64(g10) + w01*float64(g01) + w11*float64(g11))
			dst[base+2] = to01(w00*float64(bl00) + w10*float64(bl10) + w01*float64(bl01) + w11*float64(bl11))
		}
	}

	// Normalize and reorder to NCHW: for each channel c, emit row-major
	// size×size values. dst holds RGB-interleaved [0,1] floats.
	out := make([]float32, size*size*3)
	step := size * size
	for c := 0; c < 3; c++ {
		for i := 0; i < step; i++ {
			v := (dst[i*3+c] - mean[c]) / std[c]
			out[c*step+i] = float32(v)
		}
	}
	return out, nil
}

// extractPoseKeypoints produces the [1,51] Vector C input.
//
// STAND-IN: the production backbone is YOLOv8n-Pose, which outputs 17 COCO
// keypoints per detected person. Until that ONNX is dropped in, this derives
// deterministic pseudo-keypoints directly from the image: the image is split
// into 17 vertical bands (head → feet); per band it reports the normalized
// centroid x, the normalized band center y, and a confidence equal to the
// band's non-background pixel coverage. Replacing this function with the
// YOLO keypoint decoder is the ONLY change required to switch Vector C to
// real pose features.
func extractPoseKeypoints(imageBytes []byte) ([]float32, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}
	b := img.Bounds()
	if b.Dx() <= 0 || b.Dy() <= 0 {
		return nil, fmt.Errorf("image has empty bounds %v", b)
	}

	const landmarks = 17 // COCO skeleton
	w, h := b.Dx(), b.Dy()
	seen := make([]int, landmarks)
	sumX := make([]float64, landmarks)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		band := (y - b.Min.Y) * landmarks / h
		if band >= landmarks {
			band = landmarks - 1
		}
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			// Treat very dark pixels as background; this is what makes the
			// pseudo-keypoint "confidence" content-dependent.
			if (r+g+bl)/3 >= 0x0F00 {
				seen[band]++
				sumX[band] += float64(x)
			}
		}
	}

	feats := make([]float32, 0, poseKeypointDim)
	bandArea := float64(w * h / landmarks)
	if bandArea == 0 {
		bandArea = 1
	}
	for i := 0; i < landmarks; i++ {
		var cx, cy, conf float64
		if seen[i] > 0 {
			cx = sumX[i] / float64(seen[i]) / float64(w)
			cy = (float64(i) + 0.5) / float64(landmarks) // band center, normalized
			conf = float64(seen[i]) / bandArea
			if conf > 1 {
				conf = 1
			}
		}
		feats = append(feats, float32(cx), float32(cy), float32(conf))
	}
	return feats, nil
}
