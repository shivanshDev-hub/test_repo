// Package model defines the shared request/response types flowing through
// the ingestion pipeline (Phases 1-10).
package model

import "time"

// InvestigationRequest is the raw inbound payload (Phase 1 input contract).
type InvestigationRequest struct {
	AnalystID        string   `json:"analyst_id"`
	AuthorizationRef string   `json:"authorization_ref"`
	TargetType       string   `json:"target_type"`
	TargetValue      string   `json:"target_value"`
	TargetHost       string   `json:"target_host"`
	SeedTerms        []string `json:"seed_terms"`
	MaxDepth         int      `json:"max_depth"`
}

// Envelope is the Phase 1 -> Phase 2 forwarded envelope with ingress metadata.
type Envelope struct {
	ClientIP   string               `json:"client_ip"`
	AnalystID  string               `json:"analyst_id"`
	RawPayload InvestigationRequest `json:"raw_payload"`
	ReceivedAt time.Time            `json:"received_at"`
}

// RateLimitResult is the Phase 2 output.
type RateLimitResult struct {
	RateLimitStatus  string `json:"rate_limit_status"`
	AnalystID        string `json:"analyst_id"`
	RequestsInWindow int    `json:"requests_in_window"`
	RemainingQuota   int    `json:"remaining_quota"`
	RetryAfterSecs   int    `json:"retry_after,omitempty"`
}

// ValidatedRequest is the Phase 3 output (sanitized + validated).
type ValidatedRequest struct {
	InvestigationRequest
	Sanitized bool `json:"sanitized"`
}

// AuthResult is the Phase 4 output.
type AuthResult struct {
	IsAuthorized     bool   `json:"is_authorized"`
	AllowlistMatched bool   `json:"allowlist_matched"`
	Reason           string `json:"reason,omitempty"`
}

// SessionContext is the Phase 5 output - merged current + historical terms.
type SessionContext struct {
	AnalystID          string   `json:"analyst_id"`
	TargetType         string   `json:"target_type"`
	TargetHost         string   `json:"target_host"`
	MaxDepth           int      `json:"max_depth"`
	CurrentQuery       []string `json:"current_query"`
	RecentSessionQuery []string `json:"recent_session_queries"`
	SessionWindowSecs  int      `json:"session_window_seconds"`
}

// GuardDecision is the Phase 6 output.
type GuardDecision struct {
	Decision       string   `json:"decision"`
	ReasonCode     string   `json:"reason_code"`
	SanitizedSeeds []string `json:"sanitized_seeds,omitempty"`
	ViolatingTerms []string `json:"violating_terms,omitempty"`
	AuditEventID   string   `json:"audit_event_id,omitempty"`
}

// ExpansionPayload is the Phase 7 output.
type ExpansionPayload struct {
	OriginalTerm           string   `json:"original_term"`
	ExpandedTerms          []string `json:"expanded_terms"`
	PassthroughIdentifiers []string `json:"passthrough_identifiers"`
	TargetHost             string   `json:"target_host"`
	MaxDepth               int      `json:"max_depth"`
}

// VerifiedExpansion is the Phase 8 output.
type VerifiedExpansion struct {
	ExpansionID   string   `json:"expansion_id"`
	VerifiedTerms []string `json:"verified_terms"`
	StrippedTerms []string `json:"stripped_terms"`
	GuardVerdict  string   `json:"guard_verdict"`
	TargetHost    string   `json:"target_host"`
	MaxDepth      int      `json:"max_depth"`
	AnalystID     string   `json:"analyst_id"`
}

// TaskMessage is the Phase 9 -> Phase 10 stream payload.
type TaskMessage struct {
	JobID          string   `json:"job_id"`
	Target         string   `json:"target"`
	Terms          []string `json:"terms"`
	DepthRemaining int      `json:"depth_remaining"`
	EnqueuedAt     string   `json:"enqueued_at"`
}

// TaskStatus is the HSET payload for task_status:{job_id}.
type TaskStatus struct {
	Status       string `json:"status"`
	AnalystID    string `json:"analyst_id"`
	Target       string `json:"target"`
	TermsCount   string `json:"terms_count"`
	GuardVerdict string `json:"guard_verdict"`
	CreatedAt    string `json:"created_at"`
	LastUpdated  string `json:"last_updated"`
}
