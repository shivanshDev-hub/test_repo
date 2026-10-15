package sanitize

import "strings"

// safeCSSProps is the allowlist of CSS properties the render sandbox will
// accept. It is intentionally cosmetic-only: fonts, colors, spacing, layout.
// Anything capable of fetching a resource (background-image, content,
// list-style-image, cursor), navigating, or scripting is absent by design.
var safeCSSProps = map[string]bool{
	"align-content": true, "align-items": true, "align-self": true,
	"background": true, "background-attachment": true, "background-blend-mode": true,
	"background-clip": true, "background-color": true, "background-origin": true,
	"background-position": true, "background-repeat": true, "background-size": true,
	"border": true, "border-bottom": true, "border-bottom-color": true,
	"border-bottom-left-radius": true, "border-bottom-right-radius": true,
	"border-bottom-style": true, "border-bottom-width": true,
	"border-collapse": true, "border-color": true, "border-left": true,
	"border-left-color": true, "border-left-style": true, "border-left-width": true,
	"border-radius": true, "border-right": true, "border-right-color": true,
	"border-right-style": true, "border-right-width": true, "border-spacing": true,
	"border-style": true, "border-top": true, "border-top-color": true,
	"border-top-left-radius": true, "border-top-right-radius": true,
	"border-top-style": true, "border-top-width": true, "border-width": true,
	"bottom": true, "box-shadow": true, "box-sizing": true, "clear": true,
	"color": true, "column-gap": true, "direction": true, "display": true,
	"float": true, "flex": true, "flex-basis": true, "flex-direction": true,
	"flex-flow": true, "flex-grow": true, "flex-shrink": true, "flex-wrap": true,
	"font": true, "font-family": true, "font-size": true, "font-stretch": true,
	"font-style": true, "font-variant": true, "font-weight": true, "gap": true,
	"grid-column": true, "grid-column-end": true, "grid-column-gap": true,
	"grid-column-start": true, "grid-gap": true, "grid-row": true,
	"grid-row-end": true, "grid-row-gap": true, "grid-row-start": true,
	"grid-template-areas": true, "grid-template-columns": true,
	"grid-template-rows": true, "height": true, "hyphens": true,
	"justify-content": true, "left": true, "letter-spacing": true,
	"line-height": true, "list-style": true, "list-style-position": true,
	"list-style-type": true, "margin": true, "margin-bottom": true,
	"margin-left": true, "margin-right": true, "margin-top": true,
	"max-height": true, "max-width": true, "min-height": true, "min-width": true,
	"opacity": true, "order": true, "outline": true, "outline-color": true,
	"outline-offset": true, "outline-style": true, "outline-width": true,
	"overflow": true, "overflow-x": true, "overflow-y": true,
	"overflow-wrap": true, "padding": true, "padding-bottom": true,
	"padding-left": true, "padding-right": true, "padding-top": true,
	"right": true, "row-gap": true, "tab-size": true, "table-layout": true,
	"text-align": true, "text-decoration": true, "text-decoration-color": true,
	"text-decoration-line": true, "text-decoration-style": true,
	"text-indent": true, "text-overflow": true, "text-shadow": true,
	"text-transform": true, "top": true, "transform": true,
	"transform-origin": true, "transform-style": true, "unicode-bidi": true,
	"vertical-align": true, "visibility": true, "white-space": true,
	"width": true, "word-break": true, "word-spacing": true, "word-wrap": true,
	"z-index": true,
}

// cssForbiddenSubstrings are fragments that, if present in a CSS value, cause
// the whole declaration to be rejected. They cover network fetch, script
// execution, and interpreter escape hatches.
var cssForbiddenSubstrings = []string{
	"url(", "image-set(", "cross-fade(", "src(", "attr(",
	"expression(", "javascript:", "vbscript:", "behavior:", "-moz-binding",
	"@import", "@font-face", "@namespace", "link(", "progid:",
}

// CSS sanitizes a CSS declaration block (the content of a style attribute).
// It returns an empty string when nothing survives — callers then drop the
// attribute entirely.
func CSS(style string) string {
	if style == "" {
		return ""
	}
	var kept []string
	for _, decl := range strings.Split(style, ";") {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		colon := strings.Index(decl, ":")
		if colon <= 0 {
			continue
		}
		prop := strings.ToLower(strings.TrimSpace(decl[:colon]))
		value := strings.TrimSpace(decl[colon+1:])
		if value == "" {
			continue
		}
		if !safeCSSProps[prop] {
			continue
		}
		lower := strings.ToLower(value)
		blocked := false
		for _, f := range cssForbiddenSubstrings {
			if strings.Contains(lower, f) {
				blocked = true
				break
			}
		}
		if blocked {
			continue
		}
		// Keep only the safe property; the browser will ignore anything it
		// cannot parse, and we re-serialize with a quoted, escaped value.
		kept = append(kept, prop+": "+strings.ReplaceAll(value, `"`, `'`))
	}
	return strings.Join(kept, "; ")
}