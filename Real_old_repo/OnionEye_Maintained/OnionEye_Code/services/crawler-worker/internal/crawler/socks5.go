package crawler

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
)

// SocksError carries the SOCKS5 REP code so the fetch retry loop can classify
// transient Tor routing failures (general failure, network/host unreachable,
// connection refused, TTL expired) as retryable and rotate the circuit instead
// of failing the job outright. net/http wraps the dialer error, so the type
// must survive errors.As unwrapping — a plain fmt.Errorf string cannot.
type SocksError struct {
	Code byte
	Host string
}

func (e *SocksError) Error() string {
	return fmt.Sprintf("SOCKS error code=0x%02x (%s) connecting to %s", e.Code, socks5ErrorString(e.Code), e.Host)
}

// dialSOCKS5 implements a minimal SOCKS5 client that uses ATYP=0x03 (DOMAINNAME)
// to ensure .onion addresses are passed as hostnames to Tor, not resolved locally.
func dialSOCKS5(ctx context.Context, proxyAddr, targetAddr string) (net.Conn, error) {
	// Parse target host:port
	host, portStr, err := net.SplitHostPort(targetAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid target address: %w", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return nil, fmt.Errorf("invalid port: %s", portStr)
	}

	// Connect to SOCKS5 proxy
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("dial proxy: %w", err)
	}

	// SOCKS5 handshake: send version and auth methods
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	if _, err := conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
		conn.Close()
		return nil, fmt.Errorf("write handshake: %w", err)
	}

	// Read server's choice
	buf := make([]byte, 2)
	if _, err := io.ReadFull(conn, buf); err != nil {
		conn.Close()
		return nil, fmt.Errorf("read handshake response: %w", err)
	}
	if buf[0] != 0x05 {
		conn.Close()
		return nil, fmt.Errorf("unsupported SOCKS version: %d", buf[0])
	}
	if buf[1] != 0x00 {
		conn.Close()
		return nil, fmt.Errorf("auth method not accepted: %d", buf[1])
	}

	// SOCKS5 CONNECT request with ATYP=0x03 (DOMAINNAME)
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	//
	// For ATYP=0x03 (DOMAINNAME):
	//   DST.ADDR = [1 byte length][domain name]
	req := make([]byte, 0, 7+len(host))
	req = append(req, 0x05)        // VER
	req = append(req, 0x01)        // CMD: CONNECT
	req = append(req, 0x00)        // RSV
	req = append(req, 0x03)        // ATYP: DOMAINNAME
	req = append(req, byte(len(host))) // Domain length
	req = append(req, []byte(host)...)  // Domain
	req = append(req, byte(port>>8), byte(port)) // Port (big-endian)

	if _, err := conn.Write(req); err != nil {
		conn.Close()
		return nil, fmt.Errorf("write connect request: %w", err)
	}

	// Read SOCKS5 response
	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	respBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, respBuf); err != nil {
		conn.Close()
		return nil, fmt.Errorf("read connect response: %w", err)
	}
	if respBuf[0] != 0x05 {
		conn.Close()
		return nil, fmt.Errorf("invalid SOCKS version in response: %d", respBuf[0])
	}
	if respBuf[1] != 0x00 {
		conn.Close()
		return nil, &SocksError{Code: respBuf[1], Host: host}
	}

	// Read and discard BND.ADDR and BND.PORT based on ATYP
	atyp := respBuf[3]
	switch atyp {
	case 0x01: // IPv4
		io.CopyN(io.Discard, conn, 4+2)
	case 0x03: // DOMAINNAME
		lenBuf := make([]byte, 1)
		if _, err := io.ReadFull(conn, lenBuf); err != nil {
			conn.Close()
			return nil, err
		}
		io.CopyN(io.Discard, conn, int64(lenBuf[0])+2)
	case 0x04: // IPv6
		io.CopyN(io.Discard, conn, 16+2)
	default:
		conn.Close()
		return nil, fmt.Errorf("unknown address type: %d", atyp)
	}

	return conn, nil
}

func socks5ErrorString(code byte) string {
	switch code {
	case 0x01:
		return "general SOCKS server failure"
	case 0x02:
		return "connection not allowed by ruleset"
	case 0x03:
		return "network unreachable"
	case 0x04:
		return "host unreachable"
	case 0x05:
		return "connection refused"
	case 0x06:
		return "TTL expired"
	case 0x07:
		return "command not supported"
	case 0x08:
		return "address type not supported"
	default:
		return fmt.Sprintf("unknown error code: 0x%02x", code)
	}
}
