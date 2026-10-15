// Package phash computes content hashes over the RAW byte stream — before the
// bytes are ever decoded, rendered, or parsed (the "never decoded" guarantee
// from the architecture). The concrete perceptual algorithms in production are
// PDQ (images) and TMK+PDQF (video); the demo implementation uses SHA-256 over
// the raw bytes, which preserves the exact same ordering guarantee.
package phash

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

// Hasher computes a digest over raw bytes.
type Hasher interface {
	// HashOf returns the hex digest for the given raw bytes.
	HashOf(raw []byte) string
}

// SHA256Hasher is the demo exact-match hasher (raw byte stream).
type SHA256Hasher struct {
	h hash.Hash
}

// NewSHA256 returns a SHA256 hasher.
func NewSHA256() *SHA256Hasher {
	return &SHA256Hasher{h: sha256.New()}
}

// HashOf hashes raw bytes and returns the lowercase hex digest.
func (h *SHA256Hasher) HashOf(raw []byte) string {
	h.h.Reset()
	h.h.Write(raw)
	return hex.EncodeToString(h.h.Sum(nil))
}

// Zero degrades the contents of b so the ciphertext/material is not left for
// the GC to reclaim later. Used on the discard path: the buffer is zeroed the
// moment a hash match is confirmed.
func Zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}