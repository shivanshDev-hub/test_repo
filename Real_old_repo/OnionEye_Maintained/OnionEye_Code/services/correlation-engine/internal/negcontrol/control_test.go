package negcontrol_test

import (
	"testing"

	"github.com/onioneye/correlation-engine/internal/negcontrol"
)

func TestEvaluateJob_Passed(t *testing.T) {
	outcome := negcontrol.JobOutcome{
		JobID:        "job-101",
		Timeouts:     0,
		CircuitBroke: false,
		Images: []negcontrol.ImageResult{
			{ImageSrc: "img1.png", HashMatch: false, ClassifierHit: false},
		},
	}

	state := negcontrol.EvaluateJob(outcome)
	if state != negcontrol.StatePassed {
		t.Fatalf("expected state PASSED, got %s", state)
	}
}

func TestEvaluateJob_FlaggedOnHash(t *testing.T) {
	outcome := negcontrol.JobOutcome{
		JobID:        "job-102",
		Timeouts:     0,
		CircuitBroke: false,
		Images: []negcontrol.ImageResult{
			{ImageSrc: "illegal.jpg", HashMatch: true, ClassifierHit: false},
		},
	}

	state := negcontrol.EvaluateJob(outcome)
	if state != negcontrol.StateFlagged {
		t.Fatalf("expected state FLAGGED, got %s", state)
	}
}

func TestEvaluateJob_IndeterminateOnCircuitBreak(t *testing.T) {
	outcome := negcontrol.JobOutcome{
		JobID:        "job-103",
		Timeouts:     0,
		CircuitBroke: true,
		Images:       []negcontrol.ImageResult{},
	}

	state := negcontrol.EvaluateJob(outcome)
	if state != negcontrol.StateIndeterminate {
		t.Fatalf("expected state INDETERMINATE, got %s", state)
	}
}

func TestEvaluateJob_IndeterminateOnTimeouts(t *testing.T) {
	outcome := negcontrol.JobOutcome{
		JobID:        "job-104",
		Timeouts:     3,
		CircuitBroke: false,
		Images: []negcontrol.ImageResult{
			{ImageSrc: "partial.png", HashMatch: false, ClassifierHit: false},
		},
	}

	state := negcontrol.EvaluateJob(outcome)
	if state != negcontrol.StateIndeterminate {
		t.Fatalf("expected state INDETERMINATE, got %s", state)
	}
}
