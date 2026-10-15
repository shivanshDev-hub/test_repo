// Package clients provides typed HTTP clients for the internal services
// (content-shield and storage-api) that the worker depends on.
package clients

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/onioneye/crawler-worker/internal/crawler"
)

// ShieldRequest is the payload sent to content-shield POST /shield.
// For the demo the raw bytes travel in-band (base64). With a blob store,
// this would become a raw_bytes_ref handle instead.
type ShieldRequest struct {
	JobID         string `json:"job_id"`
	SourceURL     string `json:"source_url"`
	HTTPStatus    int    `json:"http_status"`
	RawBytesB64   string `json:"raw_bytes"`
	TorCircuitID  string `json:"tor_circuit_id"`
}

// ContentFlag mirrors the content_flags row written by the shield.
type ContentFlag struct {
	Category    string  `json:"category"`
	ActionTaken string  `json:"action_taken"`
	HashMatch   bool    `json:"hash_match"`
	HashAlgo    string  `json:"hash_algorithm,omitempty"`
	FlaggedAt   string  `json:"flagged_at"`
}

// ShieldResponse is the shield's verdict.
type ShieldResponse struct {
	ShieldPassed bool          `json:"shield_passed"`
	BodyRef      string        `json:"body_ref"`
	ContentFlags []ContentFlag `json:"content_flags"`
}

// ShieldClient talks to the content-shield service.
type ShieldClient struct {
	baseURL string
	http    *http.Client
}

// NewShieldClient builds a client for the content-shield service.
func NewShieldClient(baseURL string) *ShieldClient {
	return &ShieldClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Check sends raw fetched bytes to the shield. The shield runs its hash check
// BEFORE the worker parses anything, so the worker never decodes bytes the
// shield would discard.
func (s *ShieldClient) Check(ctx context.Context, f *crawler.Fetched) (*ShieldResponse, error) {
	req := ShieldRequest{
		JobID:        f.JobID,
		SourceURL:    f.URL,
		HTTPStatus:   f.HTTPStatus,
		RawBytesB64:  base64.StdEncoding.EncodeToString(f.RawBody),
		TorCircuitID: f.TorCircuitID,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := s.http.Post(s.baseURL+"/shield", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("shield request: %w", err)
	}
	defer resp.Body.Close()

	var out ShieldResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("shield response: %w", err)
	}
	return &out, nil
}
