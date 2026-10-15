package extract_test

import (
	"strings"
	"testing"

	"github.com/onioneye/crawler-worker/internal/extract"
)

// arsenalPGPKey is the exact armored PGP PUBLIC KEY BLOCK served on the
// target page arsenal6vlbdak4hz4l2ooa7643omduu6prgkfe5hjm3y5xtnjetrqad.onion.
// Its User-ID packet ("ars <ars@ars.com>") is only readable after the armor
// body is base64-decoded.
const arsenalPGPKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----

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
=D+oq
-----END PGP PUBLIC KEY BLOCK-----`

// TestExtract_ArsenalPageExpected asserts the E2E outcome the dashboard
// graph needs for scanning the arsenal page:
//   - the PGP armor must be kept as a PGPKey identifier;
//   - UID email ars@ars.com must be recovered from the decoded key body;
//   - telegram handle arsenalsoutienbot must be lifted from the t.me URL;
//   - the PGP UID alias "ars" must become a handle so two actors share an
//     identifier and a persona link can form;
//   - no fake wallets may be produced from base64 substrings inside the armor.
func TestExtract_ArsenalPageExpected(t *testing.T) {
	page := `<!DOCTYPE html>
<html>
<head><title>Arsenal - Home</title></head>
<body>
  <div class="header">
    <span class="username">arsenalsoutien</span>
  </div>
  <p>Join us on Telegram: <a href="https://t.me/arsenalsoutienbot">@arsenalsoutienbot</a></p>
  <p>Mirrors: telegram.me/arsenalsoutienbot | telegram.dog/arsenalsoutienbot</p>
  <p>PGP key:</p>
  <pre>` + arsenalPGPKey + `</pre>
  <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII1MxR4bjz9W6Ki3nCXSTA2aWJt3bJPX9+Dh47gxy8LRVFZhTB99ep1Bw8e4bUhRFZWFmZ+gq6uDgFedXWBc" alt="logo"/>
  <p>Follow policy: <span>@media only applies to CSS.</span></p>
</body>
</html>`

	ex := extract.Extract([]byte(page))

	// 1. PGP key preserved.
	if len(ex.PGPKeys) != 1 || !strings.Contains(ex.PGPKeys[0], "BEGIN PGP PUBLIC KEY BLOCK") {
		t.Fatalf("expected 1 armored PGP key, got %d", len(ex.PGPKeys))
	}

	// 2. UID email recovered from decoded armor body.
	hasUIDEmail := false
	for _, e := range ex.Emails {
		if e == "ars@ars.com" {
			hasUIDEmail = true
		}
	}
	if !hasUIDEmail {
		t.Errorf("expected PGP UID email ars@ars.com in %v", ex.Emails)
	}

	// 3. Telegram handle (t.me / telegram.me / telegram.dog variants).
	foundTG := false
	for _, h := range ex.Handles {
		if h == "arsenalsoutienbot" {
			foundTG = true
		}
	}
	if !foundTG {
		t.Errorf("expected telegram handle arsenalsoutienbot in %v", ex.Handles)
	}

	// 4. PGP UID alias becomes a handle (second actor for persona link).
	foundArs := false
	for _, h := range ex.Handles {
		if h == "ars" {
			foundArs = true
		}
	}
	if !foundArs {
		t.Errorf("expected PGP UID alias 'ars' in handles, got %v", ex.Handles)
	}

	// 5. No fake wallets from base64 inside the armor.
	if len(ex.Wallets) != 0 {
		t.Errorf("expected 0 wallets (armor base64 must be stripped), got %v", ex.Wallets)
	}
}