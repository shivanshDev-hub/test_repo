// Package expander implements Phase 7 (Qwen2.5-3B term expansion) and
// Phase 8 (output guard partial-drop sanitization).
package expander

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

// Expander calls Ollama's Qwen2.5-3B to expand seed terms (Phase 7).
type Expander struct {
	ollamaURL string
	model     string
	client    *http.Client
}

// New creates an expander.
func New(ollamaURL, model string) *Expander {
	return &Expander{ollamaURL: ollamaURL, model: model, client: &http.Client{Timeout: 30 * time.Second}}
}

// Expand generates related terms for a seed term (Phase 7).
func (e *Expander) Expand(ctx context.Context, term string, session *model.SessionContext) (*model.ExpansionPayload, error) {
	prompt := fmt.Sprintf(
		"Given the investigation seed term %q for target type %q, "+
			"generate 3-5 related search terms for dark-web OSINT investigation. "+
			"Return a JSON array of strings only. Recent session context: %s",
		term, session.TargetType, strings.Join(session.RecentSessionQuery, ", "),
	)

	body, _ := json.Marshal(map[string]any{
		"model":  e.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		e.ollamaURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	var terms []string
	if err := json.Unmarshal([]byte(out.Response), &terms); err != nil {
		// Fallback: split on commas/newlines.
		terms = strings.FieldsFunc(out.Response, func(r rune) bool {
			return r == ',' || r == '\n' || r == '"' || r == '[' || r == ']'
		})
	}

	// Passthrough identifiers (onion addresses, PGP keys, wallets) are not expanded.
	var passthrough []string
	var expanded []string
	for _, t := range terms {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if isIdentifier(t) {
			passthrough = append(passthrough, t)
		} else {
			expanded = append(expanded, t)
		}
	}

	return &model.ExpansionPayload{
		OriginalTerm:           term,
		ExpandedTerms:          expanded,
		PassthroughIdentifiers: passthrough,
		TargetHost:             session.TargetHost,
		MaxDepth:               session.MaxDepth,
	}, nil
}

// isIdentifier detects onion addresses, PGP keys, and wallet addresses.
func isIdentifier(s string) bool {
	onion := regexp.MustCompile(`^[a-z2-7]{16,56}\.onion$`)
	pgp := regexp.MustCompile(`^[A-F0-9]{16,64}$`)
	wallet := regexp.MustCompile(`^(bc1|[13])[a-zA-HJ-NP-Z0-9]{25,62}$`)
	return onion.MatchString(s) || pgp.MatchString(s) || wallet.MatchString(s)
}

// OutputGuard is the Phase 8 partial-drop sanitizer.
type OutputGuard struct {
	blockedPatterns []*regexp.Regexp
}

// NewOutputGuard creates the Phase 8 guard.
func NewOutputGuard() *OutputGuard {
	return &OutputGuard{
		blockedPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(csam|child\s*porn|child\s*abuse)\b`),
			regexp.MustCompile(`(?i)\b(exploit\s*kits?|malware\s*as\s*a\s*service)\b`),
			regexp.MustCompile(`(?i)\b(credit\s*card\s*dumps?|carding)\b`),
			regexp.MustCompile(`(?i)\b(assassination|hitman|murder\s*for\s*hire)\b`),
			regexp.MustCompile(`(?i)\b(weapons?\s*of\s*mass\s*destruction|nuclear\s*weapons?)\b`),
		},
	}
}

// Verify drops violating expanded terms and returns the verified set (Phase 8).
func (og *OutputGuard) Verify(exp *model.ExpansionPayload) *model.VerifiedExpansion {
	var verified []string
	var stripped []string

	for _, t := range exp.ExpandedTerms {
		blocked := false
		for _, re := range og.blockedPatterns {
			if re.MatchString(t) {
				stripped = append(stripped, t)
				blocked = true
				break
			}
		}
		if !blocked {
			verified = append(verified, t)
		}
	}

	// Passthrough identifiers always survive.
	verified = append(verified, exp.PassthroughIdentifiers...)

	verdict := "PASS"
	if len(stripped) > 0 {
		verdict = "PARTIAL_DROP"
	}

	return &model.VerifiedExpansion{
		ExpansionID:   fmt.Sprintf("exp_%d", time.Now().UnixNano()),
		VerifiedTerms: verified,
		StrippedTerms: stripped,
		GuardVerdict:  verdict,
		TargetHost:    exp.TargetHost,
		MaxDepth:      exp.MaxDepth,
	}
}
