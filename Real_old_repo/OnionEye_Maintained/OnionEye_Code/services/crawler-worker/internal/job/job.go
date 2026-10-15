// Package job models a crawl job as it arrives on the Redis stream.
//
// The stream message matches the contract in Architecture_Subsection_1.md:
//   XADD crawl:jobs * job_id "..." target "..." terms "[\"...\"]" depth_remaining "3" enqueued_at "..."
package job

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Job is the payload placed on the `crawl:jobs` stream by the ingestion-api.
type Job struct {
	JobID         string   `json:"job_id"`
	Target        string   `json:"target"`
	Terms         []string `json:"terms"`
	DepthRemaining int      `json:"depth_remaining"`
	EnqueuedAt    string   `json:"enqueued_at"`
}

// Parse decodes a Redis stream message's fields into a Job.
// `terms` arrives as a JSON-encoded string (json.dumps on the producer side).
func Parse(fields map[string]string) (*Job, error) {
	if fields["job_id"] == "" || fields["target"] == "" {
		return nil, errors.New("job missing required job_id/target fields")
	}

	j := &Job{
		JobID:  fields["job_id"],
		Target: fields["target"],
	}

	if raw := fields["terms"]; raw != "" {
		if err := json.Unmarshal([]byte(raw), &j.Terms); err != nil {
			return nil, fmt.Errorf("decode terms: %w", err)
		}
	}

	if d := fields["depth_remaining"]; d != "" {
		fmt.Sscanf(d, "%d", &j.DepthRemaining)
	}
	j.EnqueuedAt = fields["enqueued_at"]
	return j, nil
}
