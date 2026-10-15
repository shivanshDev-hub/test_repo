package main

import "strings"

type GraphNode struct {
	// id is what the frontend uses for graph rendering and delete requests.
	ID string `json:"id"`

	Label string `json:"label"`
	Type  string `json:"type"`

	// Optional confidence/flags used by the UI.
	Market     string   `json:"market,omitempty"`
	Confidence *float64 `json:"confidence,omitempty"`
	Flagged    bool     `json:"flagged"`

	// ActorID is the underlying actor UUID this node belongs to.
	// For type=actor nodes, it is equal to ID.
	ActorID string `json:"actor_id,omitempty"`

	// Value is the raw extracted identifier / infra indicator value
	// (email, wallet address, PGP fingerprint, onion host, IP, handle…)
	// used by the dossier OSINT JUMP actions (copy + external lookups).
	Value string `json:"value,omitempty"`

	// Kind is the raw identifiers.kind / infra_indicators.indicator_type
	// string from storage (e.g. "email", "btc_wallet", "onion_host", "ip").
	Kind string `json:"kind,omitempty"`

	// Images extracted from pages owned by this actor. The DetailPanel renders
	// them as an EVIDENCE MEDIA strip; the sidebar MEDIA tab shows them all.
	Images []GraphImage `json:"images,omitempty"`

	// Stylokey is the deterministic fixed-width (64-d) behavioral key for this
	// actor: a weighted projection of the 384-d lexical style vector fused with
	// the hour/day/tempo posting histograms (correlation-engine/pkg/stylokey).
	// Similar writers + rhythms → similar keys → the frontend seeds their graph
	// positions near each other. Empty when the actor has no behavioral profile.
	Stylokey []float64 `json:"stylokey,omitempty"`

	// Stylometry contains the high-level stylometric metrics for the entity dossier.
	Stylometry   *StylometryMeta `json:"stylometry,omitempty"`
	SampleWords  int             `json:"sample_words,omitempty"`
	Observations int             `json:"observations,omitempty"`
	Cosine       float64         `json:"cosine,omitempty"`
	Match        string          `json:"match,omitempty"`
}

type StylometryMeta struct {
	SampleWords  int     `json:"sample_words,omitempty"`
	Observations int     `json:"observations,omitempty"`
	Cosine       float64 `json:"cosine,omitempty"`
	Entropy      float64 `json:"entropy,omitempty"`
	VocabOverlap string  `json:"vocab_overlap,omitempty"`
	Match        string  `json:"match,omitempty"`
}

// anyNonZero64 reports whether at least one component of v is non-zero, i.e.
// whether a projected stylokey carries real behavioral evidence (the zero key
// means "no profile available" and must not be attached to a node).
func anyNonZero64(v []float64) bool {
	for _, x := range v {
		if x != 0 {
			return true
		}
	}
	return false
}

// GraphImage is one extracted image reference attached to a graph node
// (actor dossier). Src may be an absolute onion URL or a data URI.
type GraphImage struct {
	Src        string `json:"src"`
	Type       string `json:"type,omitempty"`
	DetectedAt string `json:"detected_at,omitempty"`
}

// nodeTypeFromKind derives the graph node type from the identifier kind stored
// in the identifiers table (e.g. "email", "wallet:BTC", "pgp_key").
func nodeTypeFromKind(kind string) string {
	k := strings.ToLower(kind)
	switch {
	case strings.Contains(k, "wallet") || strings.Contains(k, "btc") || strings.Contains(k, "xmr") || strings.Contains(k, "ltc") || strings.Contains(k, "evm"):
		return "wallet"
	case strings.Contains(k, "pgp"):
		return "pgp"
	case strings.Contains(k, "email"):
		return "email"
	case strings.Contains(k, "onion"):
		return "onion"
	case strings.Contains(k, "ip"):
		return "infra"
	case strings.Contains(k, "stylo") || strings.Contains(k, "stylometry"):
		return "stylometry"
	default:
		return "identifier"
	}
}

// displayValueFromIdentifier returns a human-readable label for a graph node
// representing an identifier, falling back to the raw value.
func displayValueFromIdentifier(kind, value string) string {
	nt := nodeTypeFromKind(kind)
	switch nt {
	case "wallet":
		currency := strings.TrimPrefix(strings.ToUpper(kind), "WALLET:")
		if currency != strings.ToUpper(kind) {
			return currency + " " + value
		}
		return value
	case "email":
		return value
	case "pgp":
		// Show first 20 chars of the key for readability.
		if len(value) > 20 {
			return value[:20] + "…"
		}
		return value
	default:
		return kind + ": " + value
	}
}

// NOTE: GraphLink uses node IDs for source/target.

type GraphLink struct {
	Source     string  `json:"source"`
	Target     string  `json:"target"`
	Type       string  `json:"type"`
	Basis      string  `json:"basis"`
	Confidence float64 `json:"confidence"`
}

type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Links []GraphLink `json:"links"`
}

type LeadResponse struct {
	LinkID                string  `json:"link_id"`
	ActorIDA              string  `json:"actor_id_a"`
	LabelA                string  `json:"label_a"`
	ActorIDB              string  `json:"actor_id_b"`
	LabelB                string  `json:"label_b"`
	Basis                 string  `json:"basis"`
	AttributionConfidence float64 `json:"attribution_confidence"`
	NegativeControlPassed bool    `json:"negative_control_passed"`
	NCState               string  `json:"nc_state"`
	CreatedAt             string  `json:"created_at"`
}

// DetailActorResponse mirrors the subset of storage-api /actors/:id used by the dashboard.
// The dashboard derives non-actor nodes conceptually by attaching actor identifiers/infra signals.
type DetailActorResponse struct {
	ActorID       string   `json:"actor_id"`
	PrimaryHandle string   `json:"primary_handle"`
	Category      *string  `json:"category"`
	Confidence    *float64 `json:"confidence,omitempty"`
	Flagged       bool     `json:"flagged"`

	Identifiers []struct {
		Kind         string `json:"kind"`
		Value        string `json:"value"`
		Type         string `json:"type"`
		DisplayValue string `json:"display_value"`
	} `json:"identifiers"`

	Infra []struct {
		IndicatorType string `json:"indicator_type"`
		OnionHost     string `json:"onion_host"`
		DisplayValue  string `json:"display_value"`
	} `json:"infra"`
}

func buildIdentifierNodeID(kind, value string) string {
	return "id_" + kind + "_" + value
}

func buildInfraNodeID(indicatorType, onionHost string) string {
	return "infra_" + indicatorType + "_" + onionHost
}
