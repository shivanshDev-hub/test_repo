// Package guard implements Phase 6: dual-tier input guard (regex + LLM).
package guard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/onioneye/ingestion-api/internal/model"
)

// Tier-1 regex patterns for high-confidence violations.
var tier1Patterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(csam|child\s*porn|child\s*abuse)\b`),
	regexp.MustCompile(`(?i)\b(exploit\s*kits?|malware\s*as\s*a\s*service)\b`),
	regexp.MustCompile(`(?i)\b(credit\s*card\s*dumps?|carding)\b`),
	regexp.MustCompile(`(?i)\b(assassination|hitman|murder\s*for\s*hire)\b`),
	regexp.MustCompile(`(?i)\b(weapons?\s*of\s*mass\s*destruction|nuclear\s*weapons?)\b`),
}

// Guard is the Phase 6 dual-tier input guard.
type Guard struct {
	ollamaURL string
	model     string
	client    *http.Client
}

// New creates a guard with the given Ollama endpoint and model.
func New(ollamaURL, model string) *Guard {
	return &Guard{
		ollamaURL: ollamaURL,
		model:     model,
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

// Check runs tier-1 regex then tier-2 LLM on the seed terms (Phase 6).
func (g *Guard) Check(ctx context.Context, sc *model.SessionContext) *model.GuardDecision {
	// Tier 1: regex scan.
	var violating []string
	var clean []string
	for _, term := range sc.CurrentQuery {
		hit := false
		for _, re := range tier1Patterns {
			if re.MatchString(term) {
				violating = append(violating, term)
				hit = true
				break
			}
		}
		if !hit {
			clean = append(clean, term)
		}
	}

	if len(violating) > 0 {
		return &model.GuardDecision{
			Decision:       "BLOCK",
			ReasonCode:     "TIER1_REGEX_VIOLATION",
			ViolatingTerms: violating,
		}
	}

	// Tier 2: LLM guard (best-effort; on failure, pass through).
	verdict, err := g.llmCheck(ctx, sc.CurrentQuery)
	if err != nil {
		return &model.GuardDecision{
			Decision:       "PASS",
			ReasonCode:     "TIER2_LLM_UNAVAILABLE",
			SanitizedSeeds: clean,
		}
	}
	if verdict == "SAFE" {
		return &model.GuardDecision{
			Decision:       "PASS",
			ReasonCode:     "TIER2_LLM_SAFE",
			SanitizedSeeds: clean,
		}
	}

	return &model.GuardDecision{
		Decision:       "BLOCK",
		ReasonCode:     "TIER2_LLM_VIOLATION",
		ViolatingTerms: sc.CurrentQuery,
	}
}

// llmCheck calls Ollama with a guard prompt.
func (g *Guard) llmCheck(ctx context.Context, terms []string) (string, error) {
	prompt := fmt.Sprintf(
		"Classify the following investigation seed terms as SAFE or UNSAFE. "+
			"UNSAFE includes CSAM, weapons, explosives, assassination, or other "+
			"imminently dangerous content. Respond with exactly SAFE or UNSAFE.\n\nTerms: %s",
		strings.Join(terms, ", "),
	)

	body, _ := json.Marshal(map[string]any{
		"model":  g.model,
		"prompt": prompt,
		"stream": false,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		g.ollamaURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	answer := strings.ToUpper(strings.TrimSpace(out.Response))
	if strings.Contains(answer, "UNSAFE") {
		return "UNSAFE", nil
	}
	return "SAFE", nil
}
