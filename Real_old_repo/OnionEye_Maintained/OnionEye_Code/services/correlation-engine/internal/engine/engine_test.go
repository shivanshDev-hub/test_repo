package engine_test

import (
	"testing"
	"time"

	"github.com/onioneye/correlation-engine/internal/engine"
)

func TestComputeComposite_HighConfidenceAlert(t *testing.T) {
	now := time.Now()
	profile := engine.ActorProfile{
		ActorID:   "actor_malicious_01",
		NCPassed:  0,
		NCFlagged: 2,
		Evidence: []engine.Evidence{
			{
				ImageSrc:        "http://darknet.onion/proof.jpg",
				ClassifierScore: 1.0,
				IsHashMatch:     true,
				FlaggedAt:       now.Add(-2 * time.Hour),
			},
		},
	}

	score, alert := engine.Evaluate(profile, engine.DefaultConfig, now)
	if !alert {
		t.Fatalf("expected alert=true for direct hash match, got alert=false with score=%f", score)
	}
	if score < 0.85 {
		t.Fatalf("expected score >= 0.85, got %f", score)
	}
}

func TestComputeComposite_FalsePositiveDampening(t *testing.T) {
	now := time.Now()
	// High number of clean crawls dampens a weak single ML flag
	profile := engine.ActorProfile{
		ActorID:   "actor_benign_operator",
		NCPassed:  120, // 120 verified clean jobs
		NCFlagged: 0,
		Evidence: []engine.Evidence{
			{
				ImageSrc:        "http://forum.onion/avatar.png",
				ClassifierScore: 0.48,
				IsHashMatch:     false,
				FlaggedAt:       now.Add(-24 * time.Hour),
			},
		},
	}

	score, alert := engine.Evaluate(profile, engine.DefaultConfig, now)
	if alert {
		t.Fatalf("expected alert=false for dampened benign actor, got alert=true with score=%f", score)
	}
	if score > 0.15 {
		t.Fatalf("expected dampened score < 0.15, got %f", score)
	}
}

func TestComputeComposite_TimeDecay(t *testing.T) {
	now := time.Now()
	// Evidence from 60 days ago with a 30-day half life should be decayed by 75%
	profile := engine.ActorProfile{
		ActorID:   "actor_old_activity",
		NCPassed:  0,
		NCFlagged: 1,
		Evidence: []engine.Evidence{
			{
				ImageSrc:        "http://oldmarket.onion/ad.jpg",
				ClassifierScore: 1.0,
				IsHashMatch:     true,
				FlaggedAt:       now.Add(-60 * 24 * time.Hour), // 2 half-lives ago
			},
		},
	}

	score := engine.ComputeComposite(profile, engine.DefaultConfig, now)
	// (0.5)^2 = 0.25
	if score < 0.23 || score > 0.27 {
		t.Fatalf("expected decayed score ~0.25, got %f", score)
	}
}
