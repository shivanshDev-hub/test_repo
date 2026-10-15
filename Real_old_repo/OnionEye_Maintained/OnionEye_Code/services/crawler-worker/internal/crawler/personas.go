// Package crawler — browser fingerprint personas.
//
// These personas replace the self-branding "Mozilla/5.0 (compatible;
// OnionEye/0.1)" UA with a rotating set of complete, browser-equivalent
// fingerprint layers: a real User-Agent plus the full header constellation
// (Accept, Sec-CH-UA*, Sec-Fetch-*) a genuine browser would send, paired with
// a uTLS ClientHello ID so the TLS layer also matches a real browser rather
// than Go's distinctive default ClientHello.
package crawler

import (
	"math/rand"

	utls "github.com/refraction-networking/utls"
)

// persona is one complete browser-alike fingerprint (request headers + TLS).
type persona struct {
	name            string
	userAgent       string
	accept          string
	acceptLanguage  string
	acceptEncoding  string
	upgradeInsecure string
	secChUA         string
	secChUAMobile   string
	secChUAPlatform string
	secFetchDest    string
	secFetchMode    string
	secFetchSite    string
	secFetchUser    string
	te              string
	hello           utls.ClientHelloID
}

// personas is the rotation pool. All are real, current browser fingerprints —
// no OnionEye-branded strings anywhere.
var personas = []persona{
	{
		name:            "chrome-120-win",
		userAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		accept:          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		acceptLanguage:  "en-US,en;q=0.9",
		acceptEncoding:  "gzip, deflate",
		upgradeInsecure: "1",
		secChUA:         `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
		secChUAMobile:   "?0",
		secChUAPlatform: `"Windows"`,
		secFetchDest:    "document",
		secFetchMode:    "navigate",
		secFetchSite:    "none",
		secFetchUser:    "?1",
		hello:           utls.HelloChrome_120,
	},
	{
		name:            "chrome-120-mac",
		userAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		accept:          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		acceptLanguage:  "en-US,en;q=0.9",
		acceptEncoding:  "gzip, deflate",
		upgradeInsecure: "1",
		secChUA:         `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
		secChUAMobile:   "?0",
		secChUAPlatform: `"macOS"`,
		secFetchDest:    "document",
		secFetchMode:    "navigate",
		secFetchSite:    "none",
		secFetchUser:    "?1",
		hello:           utls.HelloChrome_120,
	},
	{
		name:           "firefox-120-win",
		userAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
		accept:         "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		acceptLanguage: "en-US,en;q=0.5",
		acceptEncoding: "gzip, deflate",
		secFetchDest:   "document",
		secFetchMode:   "navigate",
		secFetchSite:   "none",
		te:             "trailers",
		hello:          utls.HelloFirefox_120,
	},
	{
		name:           "firefox-120-linux",
		userAgent:      "Mozilla/5.0 (X11; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
		accept:         "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		acceptLanguage: "en-US,en;q=0.5",
		acceptEncoding: "gzip, deflate",
		secFetchDest:   "document",
		secFetchMode:   "navigate",
		secFetchSite:   "none",
		te:             "trailers",
		hello:          utls.HelloFirefox_120,
	},
}

// apply stamps the persona's complete header set onto an outbound request.
func (p persona) apply(h interface{ Set(string, string) }) {
	h.Set("User-Agent", p.userAgent)
	h.Set("Accept", p.accept)
	h.Set("Accept-Language", p.acceptLanguage)
	h.Set("Accept-Encoding", p.acceptEncoding)
	h.Set("Sec-Fetch-Dest", p.secFetchDest)
	h.Set("Sec-Fetch-Mode", p.secFetchMode)
	h.Set("Sec-Fetch-Site", p.secFetchSite)
	if p.upgradeInsecure != "" {
		h.Set("Upgrade-Insecure-Requests", p.upgradeInsecure)
	}
	if p.secChUA != "" {
		h.Set("Sec-CH-UA", p.secChUA)
	}
	if p.secChUAMobile != "" {
		h.Set("Sec-CH-UA-Mobile", p.secChUAMobile)
	}
	if p.secChUAPlatform != "" {
		h.Set("Sec-CH-UA-Platform", p.secChUAPlatform)
	}
	if p.secFetchUser != "" {
		h.Set("Sec-Fetch-User", p.secFetchUser)
	}
	if p.te != "" {
		h.Set("TE", p.te)
	}
}

// randomPersona returns a uniformly random persona from the rotation pool.
func randomPersona() persona {
	return personas[rand.Intn(len(personas))]
}

// randomHelloID returns a random uTLS ClientHello fingerprint ID, mirroring a
// real browser's TLS shape instead of Go's default ClientHello.
func randomHelloID() utls.ClientHelloID {
	return personas[rand.Intn(len(personas))].hello
}