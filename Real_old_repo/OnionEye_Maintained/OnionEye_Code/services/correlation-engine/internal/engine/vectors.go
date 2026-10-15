package engine

import (
	"strings"
	"time"
)

// VectorType identifies the modality of an extracted forensic artifact.
type VectorType string

const (
	VectorEmail      VectorType = "email"
	VectorWallet     VectorType = "wallet"
	VectorPGP        VectorType = "pgp_key"
	VectorHandle     VectorType = "handle"
	VectorInfra      VectorType = "infra"
	VectorStylometry VectorType = "stylometry"
)

// ExtractedArtifact represents an extracted identifier or footprint from darknet crawls.
type ExtractedArtifact struct {
	Type        VectorType `json:"type"`
	Value       string     `json:"value"`
	Subtype     string     `json:"subtype,omitempty"`      // e.g. "BTC", "XMR", "EVM", "onion_host"
	SampleWords int        `json:"sample_words,omitempty"` // For stylometry or text blocks
	ObservedAt  time.Time  `json:"observed_at"`
}

// ScoredArtifact contains the evaluated hard score, IDF weight, time decay, and justification basis.
type ScoredArtifact struct {
	Artifact   ExtractedArtifact `json:"artifact"`
	HardScore  float64           `json:"hard_score"`  // [0.0, 1.0] Intrinsic reliability
	IDFWeight  float64           `json:"idf_weight"`  // [0.0, 1.0] Specificity / rarity penalty
	TimeDecay  float64           `json:"time_decay"`  // [0.0, 1.0] 30-day half-life exponential aging
	FinalScore float64           `json:"final_score"` // HardScore * IDFWeight * TimeDecay
	Basis      string            `json:"basis"`
}

// Known generic administrative email prefixes.
var genericEmailPrefixes = map[string]struct{}{
	"admin": {}, "administrator": {}, "support": {}, "info": {}, "contact": {},
	"abuse": {}, "webmaster": {}, "noreply": {}, "no-reply": {}, "help": {},
	"sales": {}, "billing": {}, "staff": {}, "root": {}, "moderator": {},
	"postmaster": {}, "hostmaster": {}, "team": {}, "mail": {}, "office": {},
}

// Known disposable / burner email domains.
var disposableEmailDomains = map[string]struct{}{
	"mailinator.com": {}, "tempmail.com": {}, "10minutemail.com": {},
	"guerrillamail.com": {}, "trashmail.com": {}, "yopmail.com": {},
	"dispostable.com": {}, "getairmail.com": {}, "fakemailgenerator.com": {},
	"temp-mail.org": {}, "throwawaymail.com": {}, "sharklasers.com": {},
}

// High-privacy / darknet-favored E2EE email providers.
var privacyEmailDomains = map[string]struct{}{
	"proton.me": {}, "protonmail.com": {}, "tutanota.com": {}, "tuta.com": {},
	"tuta.io": {}, "onionmail.org": {}, "secmail.pro": {}, "elude.in": {},
	"cock.li": {}, "dnmx.org": {},
}

// Known burn / genesis / faucet / test cryptocurrency addresses.
var knownGenericWallets = map[string]struct{}{
	"1a1zp1ep5qgefi2dmptftl5slmv7divfna":         {}, // Bitcoin Genesis Address
	"1111111111111111111114olvt2":                {}, // BTC Burn address
	"0x0000000000000000000000000000000000000000": {}, // EVM Zero / Burn address
	"0x000000000000000000000000000000000000dead": {}, // EVM Dead address
}

// Generic boilerplate handles with high false-positive rates.
var genericHandles = map[string]struct{}{
	"admin": {}, "administrator": {}, "mod": {}, "moderator": {}, "root": {},
	"anonymous": {}, "anon": {}, "user": {}, "guest": {}, "vendor": {},
	"support": {}, "staff": {}, "member": {}, "buyer": {}, "seller": {},
	"owner": {}, "system": {}, "bot": {},
}

// ScoreEmail intelligently scores an email based on mailbox uniqueness, domain authority, and privacy tier.
func ScoreEmail(email string) (hardScore float64, idf float64, basis string) {
	email = strings.TrimSpace(strings.ToLower(email))
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0.10, 0.10, "Malformed or invalid email structure"
	}

	mailbox := parts[0]
	domain := parts[1]

	// 1. Generic administrative prefixes (e.g. admin@, support@)
	if _, isGeneric := genericEmailPrefixes[mailbox]; isGeneric {
		return 0.15, 0.20, "Generic administrative email prefix (high false-positive risk)"
	}

	// 2. Disposable / burner domains
	if _, isDisposable := disposableEmailDomains[domain]; isDisposable {
		return 0.20, 0.25, "Disposable/burner temporary email service"
	}

	// 3. E2EE Privacy mailbox (Proton, Tuta, OnionMail)
	if _, isPrivacy := privacyEmailDomains[domain]; isPrivacy {
		if len(mailbox) >= 5 {
			return 0.95, 0.95, "Hardened E2EE privacy mailbox with specific actor username"
		}
		return 0.80, 0.85, "Hardened E2EE privacy mailbox with short username"
	}

	// 4. Custom domain or standard mailbox
	if len(mailbox) >= 5 {
		return 0.88, 0.90, "Unique personal mailbox identifier"
	}
	return 0.70, 0.75, "Short mailbox identifier"
}

// ScoreWallet intelligently scores crypto wallet addresses based on currency specificity and collision resistance.
func ScoreWallet(currency, address string) (hardScore float64, idf float64, basis string) {
	addrNorm := strings.TrimSpace(strings.ToLower(address))
	currUpper := strings.ToUpper(strings.TrimSpace(currency))

	// Check known burn/genesis/faucet addresses
	if _, isBurn := knownGenericWallets[addrNorm]; isBurn {
		return 0.05, 0.10, "Known genesis/burn or public faucet address (zero attribution value)"
	}

	switch currUpper {
	case "XMR", "MONERO":
		return 0.96, 0.98, "High-entropy Monero (XMR) stealth privacy address"
	case "BTC", "BITCOIN":
		if strings.HasPrefix(addrNorm, "bc1") {
			return 0.93, 0.92, "Native SegWit/Taproot Bech32 Bitcoin address"
		}
		return 0.88, 0.88, "Standard Base58 P2PKH/P2SH Bitcoin address"
	case "EVM", "ETH":
		return 0.89, 0.87, "Ethereum/EVM 20-byte hexadecimal wallet address"
	case "LTC", "LITECOIN":
		return 0.85, 0.85, "Litecoin (LTC) wallet address"
	default:
		return 0.75, 0.75, "Generic cryptocurrency wallet address"
	}
}

// ScorePGPKey evaluates PGP public key blocks for cryptographic strength vs demo boilerplate.
func ScorePGPKey(keyBlock string) (hardScore float64, idf float64, basis string) {
	key := strings.TrimSpace(keyBlock)
	keyLower := strings.ToLower(key)

	// Check for demo/tutorial/CTF snippets
	if strings.Contains(keyLower, "example") || strings.Contains(keyLower, "test key") || len(key) < 250 {
		return 0.10, 0.15, "Boilerplate/tutorial PGP key block or truncated snippet"
	}

	if len(key) >= 1000 && strings.Contains(key, "-----BEGIN PGP PUBLIC KEY BLOCK-----") {
		return 0.98, 0.98, "Full high-entropy cryptographic PGP public key armor block (RSA-4096/Ed25519)"
	}

	return 0.80, 0.85, "Standard PGP public key signature block"
}

// ScoreHandle evaluates handle uniqueness vs marketplace boilerplate usernames.
func ScoreHandle(handle string) (hardScore float64, idf float64, basis string) {
	h := strings.TrimSpace(strings.ToLower(handle))

	if _, isGeneric := genericHandles[h]; isGeneric {
		return 0.10, 0.15, "Generic marketplace boilerplate handle (extreme false-positive rate)"
	}

	if len(h) < 3 {
		return 0.25, 0.30, "Short ambiguous handle (< 3 characters)"
	}

	hasDigit := strings.ContainsAny(h, "0123456789")
	hasUnderscore := strings.Contains(h, "_") || strings.Contains(h, "-")
	if len(h) >= 6 && (hasDigit || hasUnderscore) {
		return 0.88, 0.90, "Distinct high-entropy adversary handle"
	}

	return 0.75, 0.80, "Standard user handle"
}

// ScoreInfra evaluates infrastructure misconfiguration paths and hidden service hosts.
func ScoreInfra(infraType, value string) (hardScore float64, idf float64, basis string) {
	v := strings.TrimSpace(strings.ToLower(value))
	t := strings.TrimSpace(strings.ToLower(infraType))

	if strings.Contains(v, "server-status") || strings.Contains(v, "phpmyadmin") || strings.Contains(v, ".git") {
		return 0.95, 0.92, "Exposed server diagnostic / configuration leak path"
	}

	if t == "onion_host" || strings.HasSuffix(v, ".onion") {
		return 0.75, 0.80, "Hidden service (.onion) destination endpoint"
	}

	return 0.65, 0.70, "General infrastructure indicator"
}

// ScoreArtifact processes an ExtractedArtifact and computes its hard score, IDF weight, and time decay.
func ScoreArtifact(artifact ExtractedArtifact, now time.Time, halfLifeDays float64) ScoredArtifact {
	if halfLifeDays <= 0 {
		halfLifeDays = 30.0
	}

	var hardScore, idf float64
	var basis string

	switch artifact.Type {
	case VectorEmail:
		hardScore, idf, basis = ScoreEmail(artifact.Value)
	case VectorWallet:
		hardScore, idf, basis = ScoreWallet(artifact.Subtype, artifact.Value)
	case VectorPGP:
		hardScore, idf, basis = ScorePGPKey(artifact.Value)
	case VectorHandle:
		hardScore, idf, basis = ScoreHandle(artifact.Value)
	case VectorInfra:
		hardScore, idf, basis = ScoreInfra(artifact.Subtype, artifact.Value)
	case VectorStylometry:
		hardScore, idf, basis = 0.62, 0.80, "Stylometric profile fingerprint (384-d lexical vector)"
	default:
		hardScore, idf, basis = 0.50, 0.50, "Generic extracted artifact"
	}

	timeDecay := 1.0
	if !artifact.ObservedAt.IsZero() {
		ageHours := now.Sub(artifact.ObservedAt).Hours()
		ageDays := ageHours / 24.0
		if ageDays > 0 {
			timeDecay = calculateDecay(ageDays, halfLifeDays)
		}
	}

	finalScore := hardScore * idf * timeDecay

	return ScoredArtifact{
		Artifact:   artifact,
		HardScore:  hardScore,
		IDFWeight:  idf,
		TimeDecay:  timeDecay,
		FinalScore: finalScore,
		Basis:      basis,
	}
}
