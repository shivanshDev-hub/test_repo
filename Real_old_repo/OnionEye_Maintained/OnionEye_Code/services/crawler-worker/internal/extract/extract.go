// Package extract implements the regex extraction map from
// Architecture_Subsection_2.md — applied AFTER the content shield passes the
// bytes. Text is normalized (zero-width strip, HTML-entity decode, bidi strip)
// before any pattern is evaluated.
package extract

import (
	"encoding/base64"
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/onioneye/crawler-worker/internal/clients"
)

// Precompiled patterns — Regex Extraction Map.
var (
	rePGP          = regexp.MustCompile(`-----BEGIN PGP (?:PUBLIC )?KEY BLOCK-----[^-]*-----END PGP (?:PUBLIC )?KEY BLOCK-----`)
	reBTC          = regexp.MustCompile(`\b(bc1[a-zA-HJ-NP-Z0-9]{39,59}|[13][a-km-zA-HJ-NP-Z1-9]{25,34})\b`)
	reEVM          = regexp.MustCompile(`\b0x[a-fA-F0-9]{40}\b`)
	reXMR          = regexp.MustCompile(`\b[48][0-9AB][1-9A-HJ-NP-Za-km-z]{93}\b`)
	reLTC          = regexp.MustCompile(`\b[LM3][a-km-zA-HJ-NP-Z1-9]{26,34}\b`)
	reOnion        = regexp.MustCompile(`\b[a-z2-7]{56}\.onion(?:/[^\s"'<>]*)?\b`)
	reEmail        = regexp.MustCompile(`\b[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}\b`)
	reURL          = regexp.MustCompile(`https?://[^\s"'<>]+`)
	reImg          = regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
	reImgDirect    = regexp.MustCompile(`\bhttps?://[^\s"'<>]+\.(?:jpg|jpeg|png|webp)\b`)
	reTs           = regexp.MustCompile(`\b\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?\b`)
	reHandleDOM    = regexp.MustCompile(`class="username"[^>]*>([^<]+)<`)
	reHandleAt     = regexp.MustCompile(`(?:^|[\s(])@([a-zA-Z0-9_\-]{3,32})\b`)
	reTelegramURL  = regexp.MustCompile(`(?i)\b(?:t\.me|telegram\.me|telegram\.dog)/([a-zA-Z0-9_]{3,32})`)
	rePGPUID       = regexp.MustCompile(`\b([a-zA-Z0-9_]{2,40})\s*<\s*([a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,})>`) // PGP User-ID "Name <email>" inside decoded body
	reDataURI      = regexp.MustCompile(`(?i)data:[^"'\s<>]*;base64,[a-z0-9+/=\s]*`)                                        // inline base64 payloads
	reTXID         = regexp.MustCompile(`\b[a-fA-F0-9]{64}\b`)
	reServerBanner = regexp.MustCompile(`(?i)\b(Server Status|Apache Status|phpinfo\(\)|X-Powered-By:|nginx/[\d\.]+)\b`)

	// Strip <script> and <style> blocks to prevent CSS at-rules and minified JS
	// from being extracted as handles, wallets, or emails.
	reScriptTag = regexp.MustCompile(`(?i)(?s)<script\b[^>]*>.*?</script>`)
	reStyleTag  = regexp.MustCompile(`(?i)(?s)<style\b[^>]*>.*?</style>`)
)

// cssAtRuleDenylist prevents CSS at-rules (e.g. @keyframes, @media) from being
// captured as user handles by the @mention regex.
var cssAtRuleDenylist = map[string]bool{
	"keyframes":           true,
	"keyframe":            true,
	"media":               true,
	"supports":            true,
	"import":              true,
	"charset":             true,
	"font-face":           true,
	"page":                true,
	"namespace":           true,
	"layer":               true,
	"property":            true,
	"container":           true,
	"scope":               true,
	"counter-style":       true,
	"font-feature-values": true,
	"color-profile":       true,
	"font-palette-values": true,
	"starting-style":      true,
	"scroll-timeline":     true,
}

// handleNoiseDenylist prevents JavaScript framework event directives (e.g. the
// Vue @click / @keydown bindings that live in HTML attributes) and boilerplate
// attribution tokens from being captured as @mention user handles.
var handleNoiseDenylist = map[string]bool{
	"click":       true,
	"dblclick":    true,
	"keydown":     true,
	"keyup":       true,
	"keypress":    true,
	"mousedown":   true,
	"mouseup":     true,
	"mousemove":   true,
	"mouseover":   true,
	"mouseenter":  true,
	"mouseleave":  true,
	"mouseout":    true,
	"scroll":      true,
	"wheel":       true,
	"resize":      true,
	"change":      true,
	"input":       true,
	"submit":      true,
	"reset":       true,
	"focus":       true,
	"blur":        true,
	"focusin":     true,
	"focusout":    true,
	"load":        true,
	"unload":      true,
	"error":       true,
	"contextmenu": true,
	"toggle":      true,
	"show":        true,
	"hide":        true,
	"mouse":       true,
	"click_":      true,
	// Common CDN / icon-library attribution tokens that frequently appear as
	// @fontawesome, @davegandy, etc. in generated HTML comments and footers.
	"fontawesome": true,
	"davegandy":   true,
	"bootstrap":   true,
	"jquery":      true,
	"webpack":     true,
	"font":        true,
	"charset":     true,
}

// emailDenylist filters placeholder / boilerplate emails that appear in HTML
// templates, documentation, and auto-generated pages.
var emailDenylist = map[string]bool{
	"you@domain.com":       true,
	"user@example.com":     true,
	"email@example.com":    true,
	"test@test.com":        true,
	"name@domain.com":      true,
	"user@domain.tld":      true,
	"your@email.com":       true,
	"john.doe@example.com": true,
	"admin@example.com":    true,
	"info@example.com":     true,
	"someone@example.com":  true,
	"no-reply@example.com": true,
	"noreply@example.com":  true,
	"admin@localhost":      true,
	"user@localhost":       true,
}

// emailDomainDenylist is a broader catch-all: any address on these domains is
// discarded regardless of the local-part.
var emailDomainDenylist = map[string]bool{
	"domain.com":  true,
	"example.com": true,
	"example.org": true,
	"example.net": true,
	"test.com":    true,
	"test.org":    true,
	"localhost":   true,
}

// Normalize applies the regex-hardening pipeline before any extraction:
// 1. Decodes HTML entities (&amp; -> &, &#x27; -> ', etc.)
// 2. Strips Zero-Width obfuscation characters (U+200B, U+200C, U+200D, U+FEFF)
// 3. Strips RTL/LTR bidi override characters (U+202A–U+202E, U+2066–U+2069)
func Normalize(raw []byte) string {
	s := string(raw)
	s = html.UnescapeString(s)
	s = stripZeroWidth(s)
	s = stripBidi(s)
	return s
}

// stripZeroWidth removes U+200B (ZWSP), U+200C (ZWNJ), U+200D (ZWJ) and U+FEFF (BOM).
func stripZeroWidth(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 0x200b, 0x200c, 0x200d, 0xfeff:
			return -1
		}
		return r
	}, s)
}

// stripBidi removes left-to-right / right-to-left override characters (U+202A–U+202E) and isolates (U+2066–U+2069).
func stripBidi(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 0x202a && r <= 0x202e) || (r >= 0x2066 && r <= 0x2069) {
			return -1
		}
		return r
	}, s)
}

// stripScriptAndStyle removes <script>...</script> and <style>...</style>
// blocks so that CSS at-rules and minified JS are never seen by extractors.
func stripScriptAndStyle(s string) string {
	s = reScriptTag.ReplaceAllString(s, "")
	s = reStyleTag.ReplaceAllString(s, "")
	return s
}

// isPlaceholderEmail reports whether the address is a boilerplate /
// documentation placeholder that should not be extracted.
func isPlaceholderEmail(email string) bool {
	if emailDenylist[strings.ToLower(email)] {
		return true
	}
	parts := strings.SplitN(strings.ToLower(email), "@", 2)
	if len(parts) == 2 && emailDomainDenylist[parts[1]] {
		return true
	}
	return false
}

// isCSSAtRule reports whether the handle matches a known CSS at-rule keyword.
func isCSSAtRule(handle string) bool {
	return cssAtRuleDenylist[strings.ToLower(handle)]
}

// isNoiseHandle reports whether the handle matches a known JavaScript framework
// event-directive or boilerplate token (e.g. Vue's @click) that is not a real
// user handle.
func isNoiseHandle(handle string) bool {
	return handleNoiseDenylist[strings.ToLower(handle)]
}

// Extract runs the full pattern set over normalized text.
func Extract(raw []byte) *clients.Extracted {
	text := Normalize(raw)
	// Strip <script> and <style> blocks before extraction to remove
	// CSS at-rules and minified JS that would produce false positives.
	text = stripScriptAndStyle(text)
	e := &clients.Extracted{}

	// 1. PGP keys (armored blocks)
	pgpBlocks := rePGP.FindAllString(text, -1)
	for _, m := range pgpBlocks {
		e.PGPKeys = append(e.PGPKeys, strings.TrimSpace(m))
	}

	// Remove armored blocks before wallet/email/handle regexes so the base64
	// payload inside the armor can never collide with wallet/email patterns.
	// Also drop inline data:...;base64, blobs (embedded images etc.) which are
	// binary payloads that could randomly form base58/hex-looking addresses.
	textNoPGP := reDataURI.ReplaceAllString(rePGP.ReplaceAllString(text, " "), " ")

	// PGP UID identities live inside the (binary) key packet — not in the
	// armored text — so decode each armor block and lift "Name <email>" UIDs.
	var pgpEmails []string
	var pgpHandleNames []string
	for _, blk := range pgpBlocks {
		emails, names := decodePGPUIDS(blk)
		pgpEmails = append(pgpEmails, emails...)
		pgpHandleNames = append(pgpHandleNames, names...)
	}

	// 2. Crypto wallets — scanned on armor-stripped text only to avoid
	//    false-positive base64 substrings.
	for _, m := range reBTC.FindAllString(textNoPGP, -1) {
		e.Wallets = append(e.Wallets, clients.Wallet{Currency: "BTC", Address: m})
	}
	for _, m := range reEVM.FindAllString(textNoPGP, -1) {
		e.Wallets = append(e.Wallets, clients.Wallet{Currency: "EVM", Address: m})
	}
	for _, m := range reXMR.FindAllString(textNoPGP, -1) {
		e.Wallets = append(e.Wallets, clients.Wallet{Currency: "XMR", Address: m})
	}
	for _, m := range reLTC.FindAllString(textNoPGP, -1) {
		e.Wallets = append(e.Wallets, clients.Wallet{Currency: "LTC", Address: m})
	}

	// 3. Infra indicators: onion hosts and server diagnostic banners
	for _, m := range reOnion.FindAllString(textNoPGP, -1) {
		e.InfraIndicators = append(e.InfraIndicators, clients.InfraIndicator{Type: "onion_host", Value: m})
	}
	for _, m := range reServerBanner.FindAllString(textNoPGP, -1) {
		e.InfraIndicators = append(e.InfraIndicators, clients.InfraIndicator{Type: "server_banner", Value: m})
	}

	// 4. Emails — plaintext body matches + decoded PGP UID emails, with
	//    placeholder / boilerplate filtering and dedupe.
	plainEmails := reEmail.FindAllString(textNoPGP, -1)
	allEmails := append(plainEmails, pgpEmails...)
	filteredEmails := make([]string, 0, len(allEmails))
	for _, m := range allEmails {
		if !isPlaceholderEmail(m) {
			filteredEmails = append(filteredEmails, m)
		}
	}
	e.Emails = dedupe(filteredEmails)

	// 5. Handles (DOM element, @mention, telegram URL path, PGP UID name) —
	//    filter CSS at-rules and framework event-directive / boilerplate noise.
	for _, m := range reHandleDOM.FindAllStringSubmatch(textNoPGP, -1) {
		if len(m) > 1 {
			h := strings.TrimSpace(m[1])
			if !isCSSAtRule(h) && !isNoiseHandle(h) {
				e.Handles = append(e.Handles, h)
			}
		}
	}
	for _, m := range reHandleAt.FindAllStringSubmatch(textNoPGP, -1) {
		if len(m) > 1 {
			h := strings.TrimSpace(m[1])
			if !isCSSAtRule(h) && !isNoiseHandle(h) {
				e.Handles = append(e.Handles, h)
			}
		}
	}
	for _, m := range reTelegramURL.FindAllStringSubmatch(textNoPGP, -1) {
		if len(m) > 1 {
			h := strings.TrimSpace(m[1])
			if !isCSSAtRule(h) && !isNoiseHandle(h) {
				e.Handles = append(e.Handles, h)
			}
		}
	}
	e.Handles = append(e.Handles, pgpHandleNames...)

	// 6. Links + images
	for _, m := range reURL.FindAllString(text, -1) {
		e.Links = append(e.Links, clients.Link{URL: m, Type: "external"})
	}
	for _, m := range reImg.FindAllStringSubmatch(text, -1) {
		if len(m) > 1 {
			e.Images = append(e.Images, clients.Image{Src: m[1], Type: "jpeg"})
		}
	}
	for _, m := range reImgDirect.FindAllString(text, -1) {
		e.Images = append(e.Images, clients.Image{Src: m, Type: "jpeg"})
	}

	// 7. Timestamps
	e.Timestamps = dedupe(reTs.FindAllString(text, -1))

	// 8. Behavioral authors — stylometric + temporal ("human") vector profiles
	//    attributed per handle. Empty when the page has no attributable prose.
	e.Authors = buildAuthorProfiles(text)

	dedupeInPlace(&e.Handles)
	dedupeInPlace(&e.PGPKeys)
	return e
}

// buildAuthorProfiles converts attributed author blocks into the wire shape
// (clients.AuthorProfile) that storage-api persists. The writing style is the
// 384-d lexical vector (reuses the pgvector-384 table), and the temporal
// family is the hour/weekday histogram plus the raw RFC3339 post times.
func buildAuthorProfiles(text string) []clients.AuthorProfile {
	blocks := AttributeAuthorBlocks(text)
	if len(blocks) == 0 {
		return nil
	}
	out := make([]clients.AuthorProfile, 0, len(blocks))
	for _, b := range blocks {
		words := CountWords(b.Text)
		times := b.Times
		ap := clients.AuthorProfile{
			Handle:      b.Handle,
			SampleWords: words,
		}
		if words > 0 {
			ap.Vector = LexicalVector(b.Text)
		}
		if len(times) > 0 {
			ap.HourHistogram = HourHistogram(times)
			ap.DayHistogram = DayHistogram(times)
			for _, t := range times {
				ap.PostTimes = append(ap.PostTimes, t.Format(time.RFC3339))
			}
		}
		out = append(out, ap)
	}
	return out
}

// decodePGPUIDS base64-decodes an armored OpenPGP block and returns any
// e-mail addresses and single-word name aliases found in its User-ID
// packets ("Name <email>"). The armored text itself only contains base64 —
// the identity lives in the decoded binary payload — so this is what keeps
// PGP-UID emails/handles from being skipped by the armor-stripping logic.
func decodePGPUIDS(armor string) (emails, names []string) {
	body := pgpArmorBody(armor)
	if body == "" {
		return nil, nil
	}
	raw, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		raw, err = base64.RawStdEncoding.DecodeString(body)
		if err != nil {
			return nil, nil
		}
	}
	for _, m := range rePGPUID.FindAllStringSubmatch(string(raw), -1) {
		if len(m) != 3 {
			continue
		}
		name, email := strings.TrimSpace(m[1]), strings.TrimSpace(m[2])
		if !isPlaceholderEmail(email) {
			emails = append(emails, email)
		}
		if isPlausibleHandle(name) {
			names = append(names, name)
		}
	}
	return dedupe(emails), dedupe(names)
}

// pgpArmorBody strips the BEGIN/END armor fence and armor header lines
// (Version:, Comment:) and the trailing CRC24 checksum line (which starts
// with "="), returning only the concatenated base64 payload.
func pgpArmorBody(block string) string {
	var b strings.Builder
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line == "" ||
			strings.HasPrefix(line, "-----") ||
			strings.HasPrefix(line, "Version:") ||
			strings.HasPrefix(line, "Comment:") ||
			strings.HasPrefix(line, "=") {
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

// isPlausibleHandle reports whether a PGP UID name is a clean single-token
// alias (letters/digits/underscores, 3–32 chars) rather than a full name
// like "John Smith" or a random base64 substring.
func isPlausibleHandle(s string) bool {
	if len(s) < 3 || len(s) > 32 {
		return false
	}
	for _, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_') {
			return false
		}
	}
	return !isCSSAtRule(s) && !isNoiseHandle(s)
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	out := in[:0]
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

func dedupeInPlace(p *[]string) { *p = dedupe(*p) }
