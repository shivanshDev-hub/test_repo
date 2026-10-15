package handler

import "testing"

// TestTermFallbackHandle verifies the fallback-handle transformation used when
// a scanned page yields no extractable handles. Fallback handles must be
// visually distinct from real handles (TERM: prefix), sanitized to a clean
// handle-like token, and reject empty/noise terms.
func TestTermFallbackHandle(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"simple term", "ShivanshSahajpal", "TERM:ShivanshSahajpal"},
		{"lowercase", "shivanshsahajpal", "TERM:shivanshsahajpal"},
		{"with space", "Ghost Market Dealer", "TERM:Ghost_Market_Dealer"},
		{"url-like input", "https://example.com/listing/42", "TERM:https_example_com_listing_42"},
		{"tabs/newlines collapse", "alpha\tbeta\ngamma", "TERM:alpha_beta_gamma"},
		{"leading/trailing junk", "  ---slim_shady---  ", "TERM:slim_shady"},
		{"empty term rejected", "", ""},
		{"all-punctuation rejected", "!!!###", ""},
		{"mixed unicode", "серый_торговец", "TERM:серый_торговец"},
		{"long term truncated", "ThisIsAnExtremelyLongInvestigationTermThatShouldBeTruncatedBecauseItIsWayTooBig", "TERM:ThisIsAnExtremelyLongInvestigationTermTh"},
		{"numeric only", "42", "TERM:42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := termFallbackHandle(tt.in); got != tt.want {
				t.Errorf("termFallbackHandle(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}