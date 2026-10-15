package validator_test

import (
"testing"

"github.com/onioneye/ingestion-api/internal/model"
"github.com/onioneye/ingestion-api/internal/validator"
)

func TestSanitize(t *testing.T) {
	// 1. Control characters & extra whitespace
	dirty := "\x00\x1f hello \t world \n \x7f"
	cleaned := validator.Sanitize(dirty)
	if cleaned != "hello world" {
		t.Fatalf("expected 'hello world', got %q", cleaned)
	}
}

func TestValidateRequest_SanitizationAndLimits(t *testing.T) {
	req := &model.InvestigationRequest{
		AnalystID:  " analyst_01 \x00",
		TargetType: "onion_domain",
		TargetValue: "http://example.onion",
		SeedTerms:  []string{" ransomware ", " \x00 ", "darknet market", "   "},
		MaxDepth:   10, // Exceeds default max of 5
	}

	validated, sanitized := validator.ValidateRequest(req, 5, 5)
	if !sanitized {
		t.Fatalf("expected sanitized=true")
	}
	if validated.AnalystID != "analyst_01" {
		t.Fatalf("expected analyst_01, got %q", validated.AnalystID)
	}
	if len(validated.SeedTerms) != 2 {
		t.Fatalf("expected 2 non-empty sanitized seed terms, got %d (%v)", len(validated.SeedTerms), validated.SeedTerms)
	}
	if validated.MaxDepth != 5 {
		t.Fatalf("expected MaxDepth clamped to 5, got %d", validated.MaxDepth)
	}
}

func TestAllowlist_Authorization(t *testing.T) {
	allowlistYAML := `
allowed_target_types:
  - onion_domain
  - clearnet_mirror
  - btc_address
`
	al, err := validator.LoadAllowlist([]byte(allowlistYAML))
	if err != nil {
		t.Fatalf("unexpected error loading allowlist: %v", err)
	}

	// Authorized
	resAuth := al.Check("onion_domain")
	if !resAuth.IsAuthorized {
		t.Fatalf("expected onion_domain to be authorized")
	}

	// Unauthorized
	resUnauth := al.Check("internal_subnet")
	if resUnauth.IsAuthorized {
		t.Fatalf("expected internal_subnet to be unauthorized")
	}
}

func TestValidateTargetValue_Charset(t *testing.T) {
	// Valid target
	if !validator.ValidateTargetValue("http://vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd.onion") {
		t.Fatalf("expected valid target value")
	}

	// Invalid target (contains forbidden shell/SQL chars)
	if validator.ValidateTargetValue("http://onion.site; rm -rf /") {
		t.Fatalf("expected invalid target with shell injection characters to fail")
	}
}
