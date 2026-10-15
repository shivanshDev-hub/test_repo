package crawler

import (
	"context"
	"log"
	"net"
	"time"

	utls "github.com/refraction-networking/utls"
)

// dialTLS dials the address over the Tor SOCKS5 endpoint and wraps the
// connection with a uTLS ClientHello that impersonates a real browser.
//
// This is the OPSEC-critical piece: without it, even with a browser-alike
// User-Agent, the TLS ClientHello emitted by Go's crypto/tls is structurally
// unique and instantly fingerprints the client as a Go program. By using uTLS
// with a real Chrome/Firefox ClientHello ID, the encrypted handshake is
// indistinguishable from the genuine browser it mimics.
//
// The HTTP transport calls this for https:// URLs (via DialTLSContext); plain
// http:// URLs keep using the raw SOCKS5 DialContext, exactly like a browser
// that never upgrades an http:// page.
func (c *Client) dialTLS(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	// 1. Raw TCP tunnel through Tor (SOCKS5 handshake).
	raw, err := dialSOCKS5(ctx, c.cfg.TorSocksAddr, addr)
	if err != nil {
		return nil, err
	}

	// 2. Rotate the browser ClientHello ID per connection.
	hello := randomHelloID()

	// 3. Wrap with uTLS and run the real TLS handshake.
	//    NextProtos is pinned to http/1.1: the Go transport does not speak
	//    HTTP/2 over a custom DialTLSContext, so advertising h2 (as a real
	//    browser would) would let the server select it and then choke on our
	//    HTTP/1.1 frames. http/1.1 keeps the handshake functional and is what
	//    the overwhelming majority of .onion endpoints negotiate anyway.
	uconn := utls.UClient(raw, &utls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
		NextProtos:         []string{"http/1.1"},
	}, hello)

	t0 := time.Now()
	if err := uconn.HandshakeContext(ctx); err != nil {
		raw.Close()
		log.Printf("[tls] uTLS handshake (persona=%s) failed for %s after %s: %v",
			hello.Str(), addr, time.Since(t0).Round(time.Millisecond), err)
		return nil, err
	}
	log.Printf("[tls] uTLS handshake OK persona=%s for %s after %s",
		hello.Str(), addr, time.Since(t0).Round(time.Millisecond))
	return uconn, nil
}