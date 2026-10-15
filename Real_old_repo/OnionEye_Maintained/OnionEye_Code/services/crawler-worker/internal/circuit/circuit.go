// Package circuit talks to the Tor control port to rotate circuits and lease
// SOCKS endpoints, matching the Tor Router Container behaviour in the
// architecture: SIGNAL NEWNYM before handing out a new circuit.
package circuit

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

// Manager wraps a Tor control-port connection.
type Manager struct {
	addr     string
	password string
}

// NewManager returns a circuit manager targeting the given control port.
// password is used for AUTHENTICATE; Tor rejects non-local control connections
// (closing the listener) unless password auth is enabled on the sidecar.
func NewManager(addr, password string) *Manager {
	return &Manager{addr: addr, password: password}
}

// NewNYM tells Tor to build fresh circuits (SIGNAL NEWNYM).
// On a fresh circuit the client should also re-allocate its SOCKS5 endpoint /
// CookieJar, since the old circuit is being discarded.
func (m *Manager) NewNYM() error {
	conn, err := net.DialTimeout("tcp", m.addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("dial control port: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Non-local control listeners require authentication; the sidecar's
	// HashedControlPassword matches this plaintext (see tor-sidecar/torrc).
	if _, err := fmt.Fprintf(conn, "AUTHENTICATE \"%s\"\r\n", m.password); err != nil {
		return err
	}
	if _, err := fmt.Fprint(conn, "SIGNAL NEWNYM\r\n"); err != nil {
		return err
	}

	reader := bufio.NewReader(conn)
	// Expect "250 OK" responses. Drain until QUIT.
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(line) // surfaced to logs; OK/err text is diagnostic
	}
	return nil
}
