package engine_test

import (
"testing"
"time"

"github.com/onioneye/correlation-engine/internal/engine"
)

func TestScoreEmail_AdminVsUnique(t *testing.T) {
	// 1. Generic admin email should receive severe false-positive penalty
	adminScore, adminIDF, adminBasis := engine.ScoreEmail("admin@alphamarket.onion")
	if adminScore > 0.20 || adminIDF > 0.25 {
		t.Fatalf("expected low score for admin email, got score=%f, idf=%f (%s)", adminScore, adminIDF, adminBasis)
	}

	// 2. Generic support email
	supportScore, _, _ := engine.ScoreEmail("support@vendorhub.onion")
	if supportScore > 0.20 {
		t.Fatalf("expected low score for support email, got score=%f", supportScore)
	}

	// 3. Disposable/burner email
	dispScore, _, dispBasis := engine.ScoreEmail("burner42@mailinator.com")
	if dispScore > 0.25 {
		t.Fatalf("expected low score for disposable email, got score=%f (%s)", dispScore, dispBasis)
	}

	// 4. Hardened E2EE privacy email with unique actor handle
	protonScore, protonIDF, protonBasis := engine.ScoreEmail("shadow_operator_99@proton.me")
	if protonScore < 0.90 || protonIDF < 0.90 {
		t.Fatalf("expected high score for unique ProtonMail, got score=%f, idf=%f (%s)", protonScore, protonIDF, protonBasis)
	}

	// 5. Tutanota E2EE email
	tutaScore, tutaIDF, _ := engine.ScoreEmail("recon_target@tuta.io")
	if tutaScore < 0.90 || tutaIDF < 0.90 {
		t.Fatalf("expected high score for unique Tutanota email, got score=%f, idf=%f", tutaScore, tutaIDF)
	}
}

func TestScoreWallet_GenesisVsTarget(t *testing.T) {
	// 1. Satoshi Genesis address (known public artifact)
	genesisScore, genesisIDF, genesisBasis := engine.ScoreWallet("BTC", "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
	if genesisScore > 0.10 || genesisIDF > 0.15 {
		t.Fatalf("expected near-zero score for Genesis BTC address, got score=%f, idf=%f (%s)", genesisScore, genesisIDF, genesisBasis)
	}

	// 2. Target Native SegWit address
	segwitScore, segwitIDF, _ := engine.ScoreWallet("BTC", "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq")
	if segwitScore < 0.90 || segwitIDF < 0.90 {
		t.Fatalf("expected high score for real target SegWit address, got score=%f, idf=%f", segwitScore, segwitIDF)
	}

	// 3. Target Monero stealth address
	xmrAddr := "44AFFq5kSiGBoZ4NMDwYtN18obc8AemS33DBLWs3H7otXMV33Wo927b43c6837130000000000000000000000000000000"
	xmrScore, xmrIDF, xmrBasis := engine.ScoreWallet("XMR", xmrAddr)
	if xmrScore < 0.95 || xmrIDF < 0.95 {
		t.Fatalf("expected top score for Monero stealth address, got score=%f, idf=%f (%s)", xmrScore, xmrIDF, xmrBasis)
	}
}

func TestScorePGPKey_DemoVsFullKey(t *testing.T) {
	// 1. Tutorial/Demo snippet
	demoKey := "-----BEGIN PGP PUBLIC KEY BLOCK-----\nVersion: Test Key\n\nexample short snippet\n-----END PGP PUBLIC KEY BLOCK-----"
	demoScore, demoIDF, demoBasis := engine.ScorePGPKey(demoKey)
	if demoScore > 0.15 || demoIDF > 0.20 {
		t.Fatalf("expected low score for demo PGP key, got score=%f, idf=%f (%s)", demoScore, demoIDF, demoBasis)
	}

	// 2. Real full armored PGP block (> 1000 chars)
	realKey := "-----BEGIN PGP PUBLIC KEY BLOCK-----\nComment: User 4096R/B4E381A\n\n"
	for len(realKey) < 1100 {
		realKey += "mQINBF2X8sEBEADJk9F10qP74gM8X5vJ4wL9p8sK7j2N0mP1q9r8s7t6u5v4w3x2y1z0\n"
	}
	realKey += "-----END PGP PUBLIC KEY BLOCK-----"

	realScore, realIDF, _ := engine.ScorePGPKey(realKey)
	if realScore < 0.95 || realIDF < 0.95 {
		t.Fatalf("expected top score for full armored PGP block, got score=%f, idf=%f", realScore, realIDF)
	}
}

func TestScoreHandle_CommonVsUnique(t *testing.T) {
	// Generic system handle
	adminScore, _, adminBasis := engine.ScoreHandle("admin")
	if adminScore > 0.15 {
		t.Fatalf("expected low score for generic 'admin' handle, got %f (%s)", adminScore, adminBasis)
	}

	// Unique alias
	uniqueScore, uniqueIDF, _ := engine.ScoreHandle("v3nd0r_shadow_88")
	if uniqueScore < 0.85 || uniqueIDF < 0.85 {
		t.Fatalf("expected high score for unique handle, got score=%f, idf=%f", uniqueScore, uniqueIDF)
	}
}

func TestEvaluateProfile_MultiVectorIntelligentCorrelation(t *testing.T) {
	now := time.Now()

	// Actor with high-entropy unique artifacts
	profile := engine.ActorProfile{
		ActorID:  "target_actor_01",
		NCPassed: 1, // Minimal baseline dampener
		Artifacts: []engine.ExtractedArtifact{
			{
				Type:       engine.VectorEmail,
				Value:      "target_ops@proton.me",
				ObservedAt: now.Add(-1 * time.Hour),
			},
			{
				Type:       engine.VectorWallet,
				Subtype:    "XMR",
				Value:      "44AFFq5kSiGBoZ4NMDwYtN18obc8AemS33DBLWs3H7otXMV33Wo927b43c6837130000000000000000000000000000000",
				ObservedAt: now.Add(-1 * time.Hour),
			},
		},
	}

	res := engine.EvaluateProfile(profile, engine.DefaultConfig, now)
	if res.CompositeScore < 0.85 {
		t.Fatalf("expected composite score >= 0.85 for high-confidence multi-vector actor, got %f", res.CompositeScore)
	}
	if !res.IsAlert {
		t.Fatalf("expected is_alert=true, got false")
	}
	if len(res.ScoredArtifacts) != 2 {
		t.Fatalf("expected 2 scored artifacts, got %d", len(res.ScoredArtifacts))
	}
}
