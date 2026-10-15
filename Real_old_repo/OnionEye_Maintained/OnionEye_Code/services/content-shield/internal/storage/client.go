// Package storage is the content-shield's client for the storage-api.
// The shield never holds DB credentials — flags are persisted via this HTTP
// boundary, matching the zero-trust architecture.
package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client talks to storage-api.
type Client struct {
	baseURL string
	http    *http.Client
}

// New builds a storage-api client.
func New(baseURL string) *Client {
	return &Client{baseURL: baseURL, http: &http.Client{Timeout: 15 * time.Second}}
}

// FlagPayload is the content_flags row shape posted to /content-flags.
type FlagPayload struct {
	JobID       string `json:"job_id"`
	SourceURL   string `json:"source_url"`
	Category    string `json:"category"`
	ActionTaken string `json:"action_taken"`
	HashMatch   bool   `json:"hash_match"`
	FlaggedAt   string `json:"flagged_at"`
}

// PersistFlag writes a flag record via storage-api (POST /content-flags).
func (c *Client) PersistFlag(ctx context.Context, p *FlagPayload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}
	resp, err := c.http.Post(c.baseURL+"/content-flags", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("storage request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("storage flag status %d", resp.StatusCode)
	}
	return nil
}