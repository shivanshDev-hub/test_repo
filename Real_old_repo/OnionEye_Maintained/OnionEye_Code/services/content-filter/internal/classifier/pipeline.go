package classifier

// pipeline.go — Stage 1→4 orchestration for the Phase 8A ML classifier.
//
// §1.3 design (verbatim from the architecture doc):
//
//	Stage 1  Hard gate (sync):   if p_minor < 0.40 → short-circuit clean pass
//	Stage 2  Parallel inference: nudity + behavior concurrently via goroutines
//	Stage 3  Composite score:    score = (0.30 × p_nudity) + (0.45 × p_minor) + (0.25 × p_exploitative)
//	Stage 4  Threshold decision:
//	           score ≥ 0.65  → content_discarded
//	           0.45 ≤ score < 0.65 → held_pending_review (review_required: true)
//	           score < 0.45  → clean pass

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Result types
// ---------------------------------------------------------------------------

// Action describes what the pipeline decided to do with the image.
type Action string

const (
	// ActionCleanPass — score < 0.45 or p_minor short-circuit; no flag row written.
	ActionCleanPass Action = "clean_pass"

	// ActionHeldPendingReview — 0.45 ≤ score < 0.65; human review required.
	ActionHeldPendingReview Action = "held_pending_review"

	// ActionContentDiscarded — score ≥ 0.65; content zeroed, flag row written.
	ActionContentDiscarded Action = "content_discarded"
)

// StageScores holds the raw per-model probabilities that went into the composite.
// These map exactly to the `stage_scores` field in the architecture's JSON output.
type StageScores struct {
	PNudity       float64 `json:"p_nudity"`
	PMinor        float64 `json:"p_minor"`
	PExploitative float64 `json:"p_exploitative"`
}

// ClassificationResult is the complete output of one pipeline run for a single image.
// The JSON tags match the content_flags shape defined in Architecture_Subsection_3 §1.1.
type ClassificationResult struct {
	// ImageSrc is the original URL/path identifier of the classified image.
	ImageSrc string `json:"image_src"`

	// ShortCircuited is true when Stage 1 gated out before Stage 2/3 ran.
	// Useful for latency metrics — majority of production calls should have this set.
	ShortCircuited bool `json:"short_circuited"`

	// Stages contains the per-model probabilities (zero for models not run).
	Stages StageScores `json:"stage_scores"`

	// CompositeScore is the Stage 3 weighted sum; 0 when short-circuited.
	CompositeScore float64 `json:"classifier_score"`

	// Action is the Stage 4 threshold decision.
	Action Action `json:"action_taken"`

	// ReviewRequired mirrors action == ActionHeldPendingReview.
	ReviewRequired bool `json:"review_required"`

	// HashMatch is always false for Phase 8A — hash matching is Phase 5's job.
	HashMatch bool `json:"hash_match"`

	// FlaggedAt records when the result was produced (ISO 8601).
	FlaggedAt time.Time `json:"flagged_at"`

	// ActorID is nullable — populated by the handler when the calling job
	// already resolved a handle; stays empty string when ambiguous.
	ActorID string `json:"actor_id,omitempty"`
}

// IsFlagged returns true for either flagged outcome (discard or review).
func (r *ClassificationResult) IsFlagged() bool {
	return r.Action == ActionContentDiscarded || r.Action == ActionHeldPendingReview
}

// ---------------------------------------------------------------------------
// Pipeline
// ---------------------------------------------------------------------------

// Pipeline holds the three model inferencers and runs the Stage 1→4 protocol.
// Zero value is not usable; use NewPipeline or set fields explicitly.
type Pipeline struct {
	// AgeModel — Vector A: InsightFace backbone → p_minor
	AgeModel Inferencer

	// NudityModel — Vector B: NudeNet v2 ONNX → p_nudity
	NudityModel Inferencer

	// BehaviorModel — Vector C: YOLOv8n-Pose + MLP → p_exploitative
	BehaviorModel Inferencer
}

// NewPipeline constructs a Pipeline from the three required inferencers.
// Returns an error if any is nil so callers don't get a panic later.
func NewPipeline(age, nudity, behavior Inferencer) (*Pipeline, error) {
	if age == nil {
		return nil, errors.New("pipeline: age model (Vector A) must not be nil")
	}
	if nudity == nil {
		return nil, errors.New("pipeline: nudity model (Vector B) must not be nil")
	}
	if behavior == nil {
		return nil, errors.New("pipeline: behavior model (Vector C) must not be nil")
	}
	return &Pipeline{
		AgeModel:      age,
		NudityModel:   nudity,
		BehaviorModel: behavior,
	}, nil
}

// Classify runs the four-stage protocol on imageBytes.
//
// Thread safety: Classify may be called concurrently from multiple goroutines.
// Each call is independent; shared state lives only inside the Inferencer
// implementations, which are required to be goroutine-safe (see Inferencer doc).
func (p *Pipeline) Classify(ctx context.Context, imageSrc string, imageBytes []byte) (*ClassificationResult, error) {
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("classify %q: empty image bytes", imageSrc)
	}

	result := &ClassificationResult{
		ImageSrc:  imageSrc,
		HashMatch: false, // Phase 8A never does hash matching — that's Phase 5
		FlaggedAt: time.Now().UTC(),
	}

	// -----------------------------------------------------------------------
	// Stage 1 — Hard gate: age model (synchronous, runs first)
	// -----------------------------------------------------------------------
	pMinor, err := p.AgeModel.Infer(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("stage1 age model (%s): %w", p.AgeModel.Name(), err)
	}
	result.Stages.PMinor = pMinor

	if pMinor < 0.40 {
		// Short-circuit: confidence the subject is an adult is high enough that
		// Stages 2-3 never run.  The overwhelming majority of production images
		// exit here, keeping inference cost near zero for clean content.
		result.ShortCircuited = true
		result.Action = ActionCleanPass
		return result, nil
	}

	// -----------------------------------------------------------------------
	// Stage 2 — Parallel inference: nudity + behavior (only when Stage 1 passed)
	// -----------------------------------------------------------------------
	type modelResult struct {
		score float64
		err   error
	}

	nudityCh := make(chan modelResult, 1)
	behaviorCh := make(chan modelResult, 1)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s, e := p.NudityModel.Infer(imageBytes)
		nudityCh <- modelResult{score: s, err: e}
	}()

	go func() {
		defer wg.Done()
		s, e := p.BehaviorModel.Infer(imageBytes)
		behaviorCh <- modelResult{score: s, err: e}
	}()

	// Wait for both goroutines even if ctx is done, so we don't leak them.
	wg.Wait()

	// Check context cancellation before reading results.
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("classify %q: context cancelled after stage 2: %w", imageSrc, err)
	}

	nudityRes := <-nudityCh
	behaviorRes := <-behaviorCh

	// Collect errors from both — report both if both failed.
	var errs []error
	if nudityRes.err != nil {
		errs = append(errs, fmt.Errorf("stage2 nudity model (%s): %w", p.NudityModel.Name(), nudityRes.err))
	}
	if behaviorRes.err != nil {
		errs = append(errs, fmt.Errorf("stage2 behavior model (%s): %w", p.BehaviorModel.Name(), behaviorRes.err))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	result.Stages.PNudity = nudityRes.score
	result.Stages.PExploitative = behaviorRes.score

	// -----------------------------------------------------------------------
	// Stage 3 — Composite score
	// Weights from §1.3:
	//   score = (0.30 × p_nudity) + (0.45 × p_minor) + (0.25 × p_exploitative)
	// -----------------------------------------------------------------------
	composite := (0.30 * result.Stages.PNudity) +
		(0.45 * result.Stages.PMinor) +
		(0.25 * result.Stages.PExploitative)
	result.CompositeScore = composite

	// -----------------------------------------------------------------------
	// Stage 4 — Threshold decision
	// -----------------------------------------------------------------------
	switch {
	case composite >= 0.65:
		// Discard path: flag as csam, content zeroed, no human review.
		result.Action = ActionContentDiscarded
		result.ReviewRequired = false

	case composite >= 0.45:
		// Borderline: held for human review.  The architecture is explicit that
		// borderline images are NEVER auto-discarded — a human must decide.
		result.Action = ActionHeldPendingReview
		result.ReviewRequired = true

	default:
		// Below 0.45: clean pass, no flag row written.
		result.Action = ActionCleanPass
	}

	return result, nil
}

// ---------------------------------------------------------------------------
// Job-level helper
// ---------------------------------------------------------------------------

// ImageInput represents one image from a parser job's `images` array.
type ImageInput struct {
	Src      string `json:"src"`
	Type     string `json:"type"` // "jpeg", "png", etc.
	RawBytes []byte `json:"-"`    // fetched by the handler before calling ClassifyJob
}

// JobRequest is the Phase 8A input shape (§1.1): one crawler job that passed
// the Phase 5 hash shield (shield_passed == true).
type JobRequest struct {
	JobID         string       `json:"job_id"`
	SourceURL     string       `json:"source_url"`
	ProvenanceTag string       `json:"provenance_tag"`
	Images        []ImageInput `json:"images"`
	HandlesOnPage []string     `json:"handles_on_page"`
}

// JobResult collects one ClassificationResult per image in the job.
type JobResult struct {
	JobID     string                  `json:"job_id"`
	SourceURL string                  `json:"source_url"`
	Results   []*ClassificationResult `json:"results"`
}

// ClassifyJob runs Classify concurrently over every image in req.
// If any single image classification returns an unrecoverable error the whole
// job fails fast.  Context cancellation propagates to all in-flight goroutines.
func (p *Pipeline) ClassifyJob(ctx context.Context, req *JobRequest) (*JobResult, error) {
	if req == nil {
		return nil, errors.New("ClassifyJob: nil job request")
	}
	if req.JobID == "" {
		req.JobID = uuid.New().String()
	}

	type indexedResult struct {
		idx int
		res *ClassificationResult
		err error
	}

	results := make([]*ClassificationResult, len(req.Images))
	ch := make(chan indexedResult, len(req.Images))

	var wg sync.WaitGroup
	for i, img := range req.Images {
		wg.Add(1)
		go func(idx int, image ImageInput) {
			defer wg.Done()
			r, e := p.Classify(ctx, image.Src, image.RawBytes)
			if r != nil && len(req.HandlesOnPage) == 1 {
				// Best-effort actor_id resolution: if exactly one handle is on
				// this page we can safely attribute the flag to that actor.
				// If zero or multiple handles are present we leave it blank per §1.7.
				r.ActorID = req.HandlesOnPage[0]
			}
			ch <- indexedResult{idx: idx, res: r, err: e}
		}(i, img)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var firstErr error
	for ir := range ch {
		if ir.err != nil && firstErr == nil {
			firstErr = ir.err
		}
		results[ir.idx] = ir.res
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return &JobResult{
		JobID:     req.JobID,
		SourceURL: req.SourceURL,
		Results:   results,
	}, nil
}
