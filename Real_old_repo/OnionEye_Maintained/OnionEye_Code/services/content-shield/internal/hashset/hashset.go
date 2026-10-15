// Package hashset holds the known-bad (CSAM/bystander) hash set.
//
// In production this is populated from a recognized hash-sharing program. The
// hackathon demo loads a small synthetic file (never real material) so the
// match path can be exercised end-to-end. The loader is a plain function so it
// can be swapped by dependency injection rather than conditional compilation.
package hashset

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// HashSet is an in-memory set of hex SHA-256 digests.
type HashSet struct {
	hashes map[string]struct{}
}

// New returns an empty hash set.
func New() *HashSet {
	return &HashSet{hashes: make(map[string]struct{})}
}

// FromFile loads a hash set from a text file, one hex digest per line.
// Blank lines and '#' comments are ignored. Duplicates are collapsed.
func FromFile(path string) (*HashSet, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open hash set: %w", err)
	}
	defer f.Close()

	hs := New()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Optional sha256: prefix is tolerated for readable corpus files.
		line = strings.TrimPrefix(line, "sha256:")
		hs.Add(line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return hs, nil
}

// Add inserts a hex digest into the set.
func (h *HashSet) Add(hexDigest string) { h.hashes[strings.ToLower(hexDigest)] = struct{}{} }

// Contains reports whether the hex digest is present.
func (h *HashSet) Contains(hexDigest string) bool {
	_, ok := h.hashes[strings.ToLower(hexDigest)]
	return ok
}

// Len returns the number of entries.
func (h *HashSet) Len() int { return len(h.hashes) }