package clients

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// StorageClient writes extracted data to the storage-api. This is the ONLY
// way the worker persists anything — it never holds DB credentials (zero-trust).
type StorageClient struct {
	baseURL string
	http    *http.Client
}

// NewStorageClient builds a client for the storage-api service.
func NewStorageClient(baseURL string) *StorageClient {
	return &StorageClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

// ScanPayload is the extraction output persisted to PostgreSQL.
type ScanPayload struct {
	JobID        string    `json:"job_id"`
	SourceURL    string    `json:"source_url"`
	Provenance   string    `json:"provenance_tag"`
	Extracted    Extracted `json:"extracted"`
	BodyRef      string    `json:"body_ref"`
	TorCircuitID string    `json:"tor_circuit_id"`
	FetchedAt    string    `json:"fetched_at"`
	Terms        []string  `json:"terms"`
}

// AuthorProfile is the behavioral fingerprint attributed to a single handle on
// a crawled page: the 384-d lexical style vector plus the temporal ("human")
// vector family (hour/weekday histograms and the raw post times that built
// them). Empty Vector with SampleWords == 0 means the page carried no
// attributable prose for that author.
type AuthorProfile struct {
	Handle        string    `json:"handle"`
	Vector        []float64 `json:"vector,omitempty"`
	SampleWords   int       `json:"sample_words"`
	HourHistogram []float64 `json:"hour_histogram,omitempty"` // 24 UTC bins (counts)
	DayHistogram  []float64 `json:"day_histogram,omitempty"`  // 7 UTC bins (counts)
	PostTimes     []string  `json:"post_times,omitempty"`     // RFC3339, attributed to this author
}

// Extracted holds the regex-extracted fields (subset of the architecture shape).
type Extracted struct {
	Handles         []string         `json:"handles"`
	PGPKeys         []string         `json:"pgp_keys"`
	Wallets         []Wallet         `json:"wallets"`
	InfraIndicators []InfraIndicator `json:"infra_indicators"`
	Emails          []string         `json:"emails"`
	Links           []Link           `json:"links"`
	Images          []Image          `json:"images"`
	Timestamps      []string         `json:"timestamps"`
	Authors         []AuthorProfile  `json:"authors,omitempty"`
}

type Wallet struct {
	Currency string `json:"currency"`
	Address  string `json:"address"`
}

type InfraIndicator struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Link struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type Image struct {
	Src  string `json:"src"`
	Type string `json:"type"`
}

// PersistScan writes a parsed scan to storage-api (POST /scans).
func (s *StorageClient) PersistScan(ctx context.Context, p *ScanPayload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return err
	}
	resp, err := s.http.Post(s.baseURL+"/scans", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("storage request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("storage persisted scan: status %d", resp.StatusCode)
	}
	return nil
}

// ScreenshotMeta is returned by storage-api when checking for an existing
// screenshot.
type ScreenshotMeta struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
}

// PersistScreenshot stores (or replaces) the PNG for a unique URL.
func (s *StorageClient) PersistScreenshot(ctx context.Context, rawURL string, png []byte) error {
	if len(png) == 0 {
		return fmt.Errorf("png is empty")
	}
	payload := map[string]string{
		"url":      rawURL,
		"png_data": jsonStringBase64(png),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/screenshots", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.http.Do(req)
	if err != nil {
		return fmt.Errorf("storage screenshot request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("storage persisted screenshot: status %d", resp.StatusCode)
	}
	return nil
}

// LookupScreenshot returns screenshot metadata for a URL (nil means 404).
func (s *StorageClient) LookupScreenshot(ctx context.Context, rawURL string) (*ScreenshotMeta, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/screenshots?url="+urlQueryEscape(rawURL), nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("storage screenshot lookup: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("storage screenshot lookup: status %d", resp.StatusCode)
	}
	var meta ScreenshotMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, err
	}
	return &meta, nil
}

func jsonStringBase64(b []byte) string {
	// avoid importing encoding/base64 at top by using stdlib via inline package.
	return base64.StdEncoding.EncodeToString(b)
}

func urlQueryEscape(s string) string {
	return url.QueryEscape(s)
}
