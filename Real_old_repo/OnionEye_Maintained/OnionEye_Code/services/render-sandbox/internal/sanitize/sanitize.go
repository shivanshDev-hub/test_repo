// Package sanitize implements the render-sandbox's zero-trust deep HTML
// sanitizer.
//
// Why this exists: a rendered page is only as safe as its input. Even after
// the content-shield's raw-byte hash gate, hostile dark-web HTML can carry
// script tags, event handlers, SVG payloads, CSS url() exfiltration probes,
// WebRTC/STUN beans, or <iframe> mounts that try to phone home.
//
// The render-sandbox therefore re-sanitizes EVERY payload it receives as
// defense-in-depth: it trusts nobody — not the upstream caller, not our own
// crawler. This pass is allowlist-based (nothing passes through by default),
// and it produces a static, network-inert, script-free DOM that is then
// rendered by a network-isolated Chromium instance.
//
// The "never decoded" guarantee from the architecture is preserved by the
// content-shield upstream. This package adds the "never re-armed" guarantee:
// even bytes that cleared the hash gate cannot re-assemble an active payload.
package sanitize

import (
"bytes"
"fmt"
"html"
"strings"

xhtml "golang.org/x/net/html"
)

// Limits — a hostile page must not be able to exhaust memory or CPU.
const (
// MaxInputBytes caps the raw HTML a caller may submit (2 MiB).
MaxInputBytes = 2 << 20
	// MaxOutputBytes caps the serialized sanitized document (4 MiB; headroom
	// is intentional for the browser chrome that wraps the body).
	MaxOutputBytes = 4 << 20
)

// ---------------------------------------------------------------------------
// Attribute policy
// ---------------------------------------------------------------------------

// safeAttr decides whether an attribute name/value pair survives. name must
// already be lowercased and must not start with "on".
func safeAttr(name, val string) (safeName, safeVal string, ok bool) {
	if !allowedAttributes[name] {
		return "", "", false
	}
	// Namespaced attributes (xlink:href, xml:base, xmlns…) can smuggle
	// activity into foreign namespaces — drop them.
	if strings.Contains(name, ":") {
		return "", "", false
	}
	if name == "width" || name == "height" || name == "border" ||
		name == "colspan" || name == "rowspan" || name == "start" {
		if !numericOnly(val) {
			return "", "", false
		}
	}
	return name, html.EscapeString(val), true
}

// numericOnly reports whether s is a plain non-negative integer.
func numericOnly(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Tree walk — the allowlist policy in action
// ---------------------------------------------------------------------------

// HTML sanitizes a full HTML document string into a render-safe document.
func HTML(raw string) (string, error) {
	if len(raw) > MaxInputBytes {
		return "", fmt.Errorf("input exceeds %d bytes", MaxInputBytes)
	}
	doc, err := xhtml.Parse(strings.NewReader(raw))
	if err != nil {
		return "", fmt.Errorf("parse html: %w", err)
	}

	clean := deepCopyNode(doc)
	stripTree(clean)
	ensureHead(clean)

	var buf bytes.Buffer
	if err := xhtml.Render(&buf, clean); err != nil {
		return "", fmt.Errorf("render sanitized html: %w", err)
	}
	out := buf.String()
	if len(out) > MaxOutputBytes {
		return "", fmt.Errorf("sanitized output exceeds %d bytes", MaxOutputBytes)
	}
	return out, nil
}

// deepCopyNode creates a full, detached clone of a node and its subtree. The
// clone has no Parent/PrevSibling/NextSibling — those links are managed by
// AppendChild/InsertBefore.
func deepCopyNode(n *xhtml.Node) *xhtml.Node {
	if n == nil {
		return nil
	}
	clone := &xhtml.Node{Type: n.Type, DataAtom: n.DataAtom, Data: n.Data, Namespace: n.Namespace}
	if len(n.Attr) > 0 {
		clone.Attr = make([]xhtml.Attribute, len(n.Attr))
		copy(clone.Attr, n.Attr)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if k := deepCopyNode(c); k != nil {
			clone.AppendChild(k)
		}
	}
	return clone
}

// stripTree walks the cloned document and removes disallowed nodes in place.
func stripTree(n *xhtml.Node) {
	if n == nil {
		return
	}
	// A dangerous element's ENTIRE subtree dies with it — never process
	// children here, otherwise script/style contents could be unwrapped or
	// re-emitted as visible text before the node itself is removed. The
	// parent's removal loop below detaches the whole subtree.
	if n.Type == xhtml.ElementNode && dangerousElements[strings.ToLower(n.Data)] {
		return
	}
	c := n.FirstChild
	for c != nil {
		next := c.NextSibling
		stripTree(c)
		c = next
	}
	// Remove dangerous elements (drop their subtrees along the way).
	c = n.FirstChild
	for c != nil {
		next := c.NextSibling
		if shouldDrop(c) {
			n.RemoveChild(c)
		}
		c = next
	}
	if n.Type == xhtml.ElementNode {
		sanitizeElementAttrs(n)
	}
}

// shouldDrop reports whether a node must be entirely removed.
func shouldDrop(n *xhtml.Node) bool {
	switch n.Type {
	case xhtml.ElementNode:
		return dangerousElements[strings.ToLower(n.Data)]
	case xhtml.CommentNode:
		// Comment content is invisible to users but a free text channel for
		// smugglers — drop it.
		return true
	// DoctypeNode for the document root is structural; text nodes are kept;
	// anything else (raw nodes, error nodes) is dropped by the caller logic.
	}
	return false
}

// sanitizeElementAttrs recomputes an element's attributes and unwraps unknown
// elements.
func sanitizeElementAttrs(n *xhtml.Node) {
tag := strings.ToLower(n.Data)
if !allowedElements[tag] {
unwrapElement(n)
return
}
out := n.Attr[:0]
for _, a := range n.Attr {
name := strings.ToLower(a.Key)
if strings.HasPrefix(name, "on") {
continue
}
if name == "style" {
if css := CSS(a.Val); css != "" {
out = append(out, xhtml.Attribute{Key: "style", Val: css})
}
continue
}
safeName, safeVal, ok := safeAttr(name, a.Val)
if !ok {
continue
}
out = append(out, xhtml.Attribute{Key: safeName, Val: safeVal})
}
n.Attr = out
}

// unwrapElement replaces an element node with its children, preserving the
// children's content while discarding the non-allowlisted container.
func unwrapElement(n *xhtml.Node) {
	if n.Parent == nil {
		return
	}
	parent := n.Parent
	c := n.FirstChild
	for c != nil {
		next := c.NextSibling
		// RemoveChild detaches c from n, making it a free node that can be
		// re-parented by InsertBefore.
		n.RemoveChild(c)
		parent.InsertBefore(c, n)
		c = next
	}
	parent.RemoveChild(n)
}

// ensureHead guarantees the document has an <html><head> skeleton so the
// renderer can inject its hardening CSP meta without ambiguity.
func ensureHead(root *xhtml.Node) {
	var htmlRoot *xhtml.Node
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == xhtml.ElementNode && c.Data == "html" {
			htmlRoot = c
			break
		}
	}
	if htmlRoot == nil {
		htmlRoot = &xhtml.Node{Type: xhtml.ElementNode, Data: "html"}
		root.AppendChild(htmlRoot)
	}
	for c := htmlRoot.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == xhtml.ElementNode && c.Data == "head" {
			return
		}
	}
	headNode := &xhtml.Node{Type: xhtml.ElementNode, Data: "head"}
	htmlRoot.InsertBefore(headNode, htmlRoot.FirstChild)
}
