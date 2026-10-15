// Package validator implements Phase 3 (sanitize+validate) and Phase 4 (allowlist auth).
package validator

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/onioneye/ingestion-api/internal/model"
)

var (
	reControl    = regexp.MustCompile(`[\x00-\x1f\x7f]`)
	reWhitespace = regexp.MustCompile(`\s+`)
	reTarget     = regexp.MustCompile(`^[a-zA-Z0-9._:@/-]{1,256}$`)
)

// Sanitize strips control chars and collapses whitespace in a term.
func Sanitize(s string) string {
	s = reControl.ReplaceAllString(s, "")
	s = reWhitespace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// ValidateRequest sanitizes and validates the raw payload (Phase 3).
// Returns the validated request and whether any sanitization occurred.
func ValidateRequest(req *model.InvestigationRequest, maxSeedTerms, maxDepthMax int) (*model.ValidatedRequest, bool) {
	sanitized := false
	out := &model.ValidatedRequest{}

	out.AnalystID = Sanitize(req.AnalystID)
	out.AuthorizationRef = Sanitize(req.AuthorizationRef)
	out.TargetType = Sanitize(req.TargetType)
	out.TargetValue = Sanitize(req.TargetValue)
	out.TargetHost = Sanitize(req.TargetHost)

	if out.AnalystID != req.AnalystID || out.TargetType != req.TargetType ||
		out.TargetValue != req.TargetValue || out.TargetHost != req.TargetHost {
		sanitized = true
	}

	// Seed terms: sanitize each, drop empties, cap at max.
	for _, t := range req.SeedTerms {
		clean := Sanitize(t)
		if clean == "" {
			sanitized = true
			continue
		}
		if len(out.SeedTerms) >= maxSeedTerms {
			sanitized = true
			break
		}
		out.SeedTerms = append(out.SeedTerms, clean)
		if clean != t {
			sanitized = true
		}
	}

	// Depth clamping.
	out.MaxDepth = req.MaxDepth
	if out.MaxDepth <= 0 {
		out.MaxDepth = 3
		sanitized = true
	}
	if out.MaxDepth > maxDepthMax {
		out.MaxDepth = maxDepthMax
		sanitized = true
	}

	out.Sanitized = sanitized
	return out, sanitized
}

// Allowlist is the Phase 4 authorization check against target_allowlist.yaml.
type Allowlist struct {
	AllowedTargetTypes map[string]bool
}

// LoadAllowlist parses the YAML allowlist file.
func LoadAllowlist(data []byte) (*Allowlist, error) {
	var doc struct {
		AllowedTargetTypes []string `yaml:"allowed_target_types"`
	}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	a := &Allowlist{AllowedTargetTypes: make(map[string]bool)}
	for _, t := range doc.AllowedTargetTypes {
		a.AllowedTargetTypes[t] = true
	}
	return a, nil
}

// Check verifies the target type is allowlisted (Phase 4).
func (a *Allowlist) Check(targetType string) *model.AuthResult {
	if a.AllowedTargetTypes[targetType] {
		return &model.AuthResult{IsAuthorized: true, AllowlistMatched: true}
	}
	return &model.AuthResult{IsAuthorized: false, AllowlistMatched: false, Reason: "target_type not in allowlist"}
}

// ValidateTargetValue enforces the target value charset (Phase 4 hardening).
func ValidateTargetValue(v string) bool {
	return reTarget.MatchString(v)
}

// RedisClient is a minimal interface for session/stream deps.
type RedisClient interface {
	redis.Cmdable
}

// GinContext is a type alias for handler signatures.
type GinContext = gin.Context
