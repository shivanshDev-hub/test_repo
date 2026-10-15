package negcontrol

// control.go implements Phase 9 Engine 1: the Negative Control State Machine.
// From Architecture §3.1: Engine 1 identifies negative control candidates by
// evaluating the terminal state of parsing jobs.

type EvaluationState string

const (
	// StatePassed means the job completed fully and returned 0 abusive content flags.
	// These jobs serve as the denominator (NC_passed) for Engine 2's signal dampening.
	StatePassed EvaluationState = "PASSED"

	// StateFlagged means at least one item in the job hit a CSAM hash or ML classifier.
	// Actors with FLAGGED > 0 can never be used as a clean baseline.
	StateFlagged EvaluationState = "FLAGGED"

	// StateIndeterminate means the job failed mid-flight (circuit breaker, 404, auth wall).
	// We cannot boldly claim it is PASSED, so it is ignored for NC baseline calculations.
	StateIndeterminate EvaluationState = "INDETERMINATE"
)

// ImageResult represents the Phase 5/Phase 8A outcome of a single image.
type ImageResult struct {
	ImageSrc       string
	HashMatch      bool
	ClassifierHit  bool // score >= 0.45
	ShortCircuited bool // p_minor < 0.40
}

// JobOutcome represents the execution result of the crawler.
type JobOutcome struct {
	JobID        string
	Timeouts     int
	CircuitBroke bool
	Images       []ImageResult
}

// EvaluateJob applies the Engine 1 state machine to a completed job.
func EvaluateJob(outcome JobOutcome) EvaluationState {
	if outcome.CircuitBroke {
		return StateIndeterminate
	}

	for _, img := range outcome.Images {
		if img.HashMatch || img.ClassifierHit {
			return StateFlagged
		}
	}

	// If there were network timeouts but no flags, we can't definitively say
	// the *entire* page was clean, because we failed to parse some of it.
	if outcome.Timeouts > 0 {
		return StateIndeterminate
	}

	// Job fully completed, no timeouts, no flags.
	return StatePassed
}
