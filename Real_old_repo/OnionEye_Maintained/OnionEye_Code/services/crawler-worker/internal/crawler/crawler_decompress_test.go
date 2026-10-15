package crawler

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"testing"
)

func TestDecodeContentEncodingGzip(t *testing.T) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write([]byte("<html>arsenal6vlbdak4hz4l2ooa7643omduu6prgkfe5hjm3y5xtnjetrqad.onion</html>")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	out, err := decodeContentEncoding("gzip", buf.Bytes(), 1<<20)
	if err != nil {
		t.Fatalf("decode gzip: %v", err)
	}
	if !bytes.Contains(out, []byte("arsenal6")) {
		t.Fatalf("expected decoded HTML, got %q", out)
	}
}

func TestDecodeContentEncodingDeflateZlib(t *testing.T) {
	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write([]byte("wallets=1 links=3")); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := decodeContentEncoding("deflate", buf.Bytes(), 1<<20)
	if err != nil {
		t.Fatalf("decode deflate(zlib): %v", err)
	}
	if string(out) != "wallets=1 links=3" {
		t.Fatalf("unexpected decode: %q", out)
	}
}

func TestDecodeContentEncodingDeflateRaw(t *testing.T) {
	var buf bytes.Buffer
	fw, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fw.Write([]byte("raw deflate payload")); err != nil {
		t.Fatal(err)
	}
	if err := fw.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := decodeContentEncoding("deflate", buf.Bytes(), 1<<20)
	if err != nil {
		t.Fatalf("decode deflate(raw): %v", err)
	}
	if string(out) != "raw deflate payload" {
		t.Fatalf("unexpected decode: %q", out)
	}
}

func TestDecodeContentEncodingUnsupportedKeepsRaw(t *testing.T) {
	if _, err := decodeContentEncoding("br", []byte("anything"), 1<<20); err == nil {
		t.Fatal("expected error for unsupported encoding")
	}
}

func TestDecodeContentEncodingCapsOutputSize(t *testing.T) {
	big := bytes.Repeat([]byte("a"), 1000)
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	if _, err := zw.Write(big); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := decodeContentEncoding("gzip", buf.Bytes(), 100)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) > 101 {
		t.Fatalf("output not capped: len=%d", len(out))
	}
}

func TestDecodeContentEncodingCorruptGzipFails(t *testing.T) {
	if _, err := decodeContentEncoding("gzip", []byte{0x1f, 0x8b, 0x08, 0x00, 0xff}, 1<<20); err == nil {
		t.Fatal("expected error for corrupt gzip")
	}
}