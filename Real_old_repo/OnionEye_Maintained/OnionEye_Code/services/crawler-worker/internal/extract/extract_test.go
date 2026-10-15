package extract_test

import (
"strings"
"testing"

"github.com/onioneye/crawler-worker/internal/extract"
)

func TestNormalize_AntiObfuscation(t *testing.T) {
	// 1. Zero-width spaces & BOM stripping
	obfuscatedEmail := "target\u200B_operator\u200C@\u200Dproton\uFEFF.me"
	normalized := extract.Normalize([]byte(obfuscatedEmail))
	if normalized != "target_operator@proton.me" {
		t.Fatalf("expected stripped zero-width characters, got %q", normalized)
	}

	// 2. HTML Entity decoding
	htmlEntities := "contact&amp;support&#x40;domain.com &lt;tag&gt;"
	normalizedHTML := extract.Normalize([]byte(htmlEntities))
	if !strings.Contains(normalizedHTML, "contact&support@domain.com") {
		t.Fatalf("expected decoded HTML entities, got %q", normalizedHTML)
	}

	// 3. BiDi / LTR / RTL override character stripping
	bidiText := "safe\u202Ereversed\u202Ctext"
	normalizedBidi := extract.Normalize([]byte(bidiText))
	if strings.ContainsRune(normalizedBidi, 0x202e) || strings.ContainsRune(normalizedBidi, 0x202c) {
		t.Fatalf("expected stripped bidi override characters, got %q", normalizedBidi)
	}
}

func TestExtract_AllForensicVectors(t *testing.T) {
	sampleHTML := `
<!DOCTYPE html>
<html>
<head><title>Marketplace Vendor Profile - nginx/1.24.0</title></head>
<body>
  <div class="header">
    <span class="username">dark_vendor_99</span>
    <span>Posted on: 2026-09-06T14:32:00Z</span>
  </div>

  <p>Contact me at <b>vendor_ops&#x40;proton.me</b> or reach @shadow_partner for bulk.</p>
  <p>Mirror: http://vww6ybal4bd7szmgncyruucpgfkqahzddi37ktceo3ah7ngmcopnpyyd.onion/store</p>

  <h3>Payment Methods</h3>
  <ul>
    <li>BTC Native SegWit: bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq</li>
    <li>BTC Legacy: 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa</li>
    <li>Monero XMR: 48edfHu7V9Z84YzzMa6fUueoXo6A5619bnjaWofzsmmfZtEwKkgfmgeiP15817o8ogjy5QnoAZgUCRpWPrEA6W6A1ZeUBPr</li>
    <li>Ethereum EVM: 0x71C7656EC7ab88b098defB751B7401B5f6d8976F</li>
    <li>Litecoin LTC: LQTpS3VaYTuoWdbWcuLn7t7y6xN3JpGjT1</li>
  </ul>

  <div class="pgp">
-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: Keybase OpenPGP v2.1.13
Comment: https://keybase.io/crypto

mQINBF2X8sEBEADJk9F10qP74gM8X5vJ4wL9p8sK7j2N0mP1q9r8s7t6u5v4w3x2
-----END PGP PUBLIC KEY BLOCK-----
  </div>

  <img src="http://market.onion/images/proof1.png" alt="proof" />
  <p>Direct image link: https://imagehost.onion/uploads/cert2.jpg</p>
  <p>Server Status page enabled at origin.</p>
</body>
</html>
`

	ex := extract.Extract([]byte(sampleHTML))

	// 1. Wallets
	if len(ex.Wallets) < 5 {
		t.Fatalf("expected at least 5 wallets, got %d (%+v)", len(ex.Wallets), ex.Wallets)
	}
	foundBTC, foundXMR, foundEVM, foundLTC := false, false, false, false
	for _, w := range ex.Wallets {
		switch w.Currency {
		case "BTC":
			foundBTC = true
		case "XMR":
			foundXMR = true
		case "EVM":
			foundEVM = true
		case "LTC":
			foundLTC = true
		}
	}
	if !foundBTC || !foundXMR || !foundEVM || !foundLTC {
		t.Fatalf("missing expected cryptocurrency types: BTC=%v XMR=%v EVM=%v LTC=%v", foundBTC, foundXMR, foundEVM, foundLTC)
	}

	// 2. PGP Key
	if len(ex.PGPKeys) == 0 || !strings.Contains(ex.PGPKeys[0], "BEGIN PGP PUBLIC KEY BLOCK") {
		t.Fatalf("expected extracted PGP key, got %v", ex.PGPKeys)
	}

	// 3. Email
	if len(ex.Emails) == 0 || ex.Emails[0] != "vendor_ops@proton.me" {
		t.Fatalf("expected extracted email vendor_ops@proton.me, got %v", ex.Emails)
	}

	// 4. Handles
	foundDOM, foundAt := false, false
	for _, h := range ex.Handles {
		if h == "dark_vendor_99" {
			foundDOM = true
		}
		if h == "shadow_partner" {
			foundAt = true
		}
	}
	if !foundDOM || !foundAt {
		t.Fatalf("expected extracted handles dark_vendor_99 and shadow_partner, got %v", ex.Handles)
	}

	// 5. Infra Indicators (onion host + server banner)
	if len(ex.InfraIndicators) < 2 {
		t.Fatalf("expected at least 2 infra indicators, got %d", len(ex.InfraIndicators))
	}

	// 6. Timestamps
	if len(ex.Timestamps) == 0 || !strings.Contains(ex.Timestamps[0], "2026-09-06") {
		t.Fatalf("expected extracted ISO 8601 timestamp, got %v", ex.Timestamps)
	}

	// 7. Images
	if len(ex.Images) < 2 {
		t.Fatalf("expected at least 2 images, got %d", len(ex.Images))
	}
}

// TestExtract_StripScriptStyle verifies that <script> and <style> blocks are
// removed before extraction, preventing CSS at-rules from becoming handles and
// minified JS from generating false-positive wallets / emails.
func TestExtract_StripScriptStyle(t *testing.T) {
	html := `
<html>
<head>
  <style>
    @keyframes fade { from { opacity: 0; } to { opacity: 1; } }
    @media (max-width: 768px) { body { font-size: 14px; } }
    .user { color: red; }
  </style>
</head>
<body>
  <span class="username">real_user_alice</span>
  <p>Contact: alice@proton.me</p>
  <script>
    var config = { api_key: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" };
    fetch("/api/data").then(function(r) { return r.json(); });
  </script>
  <p>My BTC: bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq</p>
</body>
</html>
`

	ex := extract.Extract([]byte(html))

	// Should NOT contain CSS at-rules as handles.
	for _, h := range ex.Handles {
		if h == "keyframes" || h == "media" {
			t.Errorf("CSS at-rule %q should not appear in handles", h)
		}
	}

	// Should contain the real handle.
	found := false
	for _, h := range ex.Handles {
		if h == "real_user_alice" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected handle 'real_user_alice' in %v", ex.Handles)
	}

	// Should contain the real email but NOT any from <script>.
	for _, e := range ex.Emails {
		if e == "alice@proton.me" {
			continue // expected
		}
		t.Errorf("unexpected email %q (likely from stripped script block)", e)
	}

	// Should contain the real BTC address from the page body, NOT the one
	// embedded in the JS config object.
	for _, w := range ex.Wallets {
		if w.Currency == "BTC" && w.Address == "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq" {
			return // correct — found the real address
		}
	}
	t.Errorf("expected real BTC wallet in %v", ex.Wallets)
}

// TestExtract_FilterFrameworkNoise verifies that JavaScript framework event
// directives (Vue's @click / @keydown / @mousedown) and boilerplate CDN
// attribution tokens (@fontawesome, @davegandy) are NOT captured as handles,
// while genuine @mention user handles are still extracted.
func TestExtract_FilterFrameworkNoise(t *testing.T) {
	html := `
<html>
<body>
  <button @click="submitOrder">Place Order</button>
  <input @keydown.enter="search" @mousedown="select" @submit.prevent="send" />
  <!-- @fontawesome and @davegandy attribution from the CDN credit -->
  <script>
    // bootstrap and jquery tokens should also be ignored
  </script>
  <p>Reach the operator directly at @real_dealer or @wired_contact for bulk.</p>
  <span class="username">legit_member</span>
</body>
</html>
`
	ex := extract.Extract([]byte(html))

	noiseTokens := map[string]bool{
		"click": true, "keydown": true, "mousedown": true, "submit": true,
		"fontawesome": true, "davegandy": true, "bootstrap": true, "jquery": true,
	}
	for _, h := range ex.Handles {
		if noiseTokens[h] {
			t.Errorf("framework/boilerplate token %q should NOT appear in handles: %v", h, ex.Handles)
		}
	}

	// Genuine handles must still be present.
	seen := map[string]bool{}
	for _, h := range ex.Handles {
		seen[h] = true
	}
	for _, want := range []string{"real_dealer", "wired_contact", "legit_member"} {
		if !seen[want] {
			t.Errorf("expected genuine handle %q in %v", want, ex.Handles)
		}
	}
}

// TestExtract_FilterPlaceholderEmails verifies that boilerplate / template
// emails are discarded while real addresses are kept.
func TestExtract_FilterPlaceholderEmails(t *testing.T) {
	text := `
		Real contact: bob@darkmarket.onion
		Templates: you@domain.com, user@example.com, admin@localhost
		More: info@example.org, noreply@test.com
	`
	ex := extract.Extract([]byte(text))

	for _, e := range ex.Emails {
		if e == "you@domain.com" || e == "user@example.com" || e == "admin@localhost" || e == "info@example.org" || e == "noreply@test.com" {
			t.Errorf("placeholder email %q should have been filtered", e)
		}
	}
	found := false
	for _, e := range ex.Emails {
		if e == "bob@darkmarket.onion" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected real email bob@darkmarket.onion in %v", ex.Emails)
	}
}
