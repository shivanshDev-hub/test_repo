package guard_test

import (
"context"
"testing"

"github.com/onioneye/ingestion-api/internal/guard"
"github.com/onioneye/ingestion-api/internal/model"
)

func TestGuard_Tier1RegexBlock(t *testing.T) {
	// Guard instance pointing to dummy Ollama URL (should not be reached because Tier 1 blocks immediately)
	g := guard.New("http://localhost:11434", "llama-guard3:1b")

	testCases := []struct {
		name     string
		term     string
		violates bool
	}{
		{"CSAM violation", "underage child abuse material", true},
		{"Exploit kit", "blackhole exploit kit download", true},
		{"Carding", "fresh credit card dumps with cvv", true},
		{"Hitman / Assassination", "darknet hitman murder for hire", true},
		{"WMD / Nuclear", "smuggled nuclear weapons", true},
		{"Legitimate OSINT query", "ransomware payment wallet btc", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
sc := &model.SessionContext{
				CurrentQuery: []string{tc.term},
			}
			decision := g.Check(context.Background(), sc)
			if tc.violates {
				if decision.Decision != "BLOCK" || decision.ReasonCode != "TIER1_REGEX_VIOLATION" {
					t.Fatalf("expected TIER1_REGEX_VIOLATION BLOCK for %q, got %+v", tc.term, decision)
				}
			} else {
				// Non-violating term should pass Tier 1 (it will hit Tier 2 LLM and fail-open with TIER2_LLM_UNAVAILABLE if Ollama is offline)
				if decision.Decision == "BLOCK" && decision.ReasonCode == "TIER1_REGEX_VIOLATION" {
					t.Fatalf("unexpected Tier 1 block for legitimate term %q", tc.term)
				}
			}
		})
	}
}
