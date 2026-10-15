package sanitize

// ---------------------------------------------------------------------------
// Element policy — nothing is allowed by default.
// ---------------------------------------------------------------------------

// allowedElements are cosmetic layout/presentation elements that cannot
// execute, fetch, or navigate. Everything else is either dropped entirely
// (dangerous set) or transparently unwrapped (unknown elements — their text
// children are kept, the container is removed).
var allowedElements = map[string]bool{
	"a": true, "abbr": true, "address": true, "article": true, "aside": true,
	"b": true, "bdi": true, "bdo": true, "blockquote": true, "body": true, "br": true,
	"caption": true, "cite": true, "code": true, "col": true, "colgroup": true,
	"data": true, "dd": true, "del": true, "details": true, "dfn": true,
	"div": true, "dl": true, "dt": true, "em": true, "figcaption": true,
	"figure": true, "footer": true, "h1": true, "h2": true, "h3": true,
	"h4": true, "h5": true, "h6": true, "head": true, "header": true, "hgroup": true, "html": true,
	"hr": true, "i": true, "ins": true, "kbd": true, "li": true, "main": true,
	"mark": true, "nav": true, "ol": true, "p": true, "pre": true, "q": true,
	"rp": true, "rt": true, "ruby": true, "s": true, "samp": true,
	"section": true, "small": true, "span": true, "strong": true, "sub": true,
	"summary": true, "sup": true, "table": true, "tbody": true, "td": true,
	"tfoot": true, "th": true, "thead": true, "time": true, "tr": true,
	"tt": true, "u": true, "ul": true, "var": true, "wbr": true,
}

// dangerousElements are dropped WITH their entire subtree. A subtree drop is
// important: dangling text from a <script> is meaningless, and <style>
// contents must never be re-emitted into the page.
var dangerousElements = map[string]bool{
	"applet": true, "audio": true, "base": true, "canvas": true, "dialog": true,
	"embed": true, "form": true, "frame": true, "frameset": true, "iframe": true,
	"image": true, "input": true, "keygen": true, "label": true, "link": true,
	"marquee": true, "math": true, "meta": true, "meter": true, "noembed": true,
	"noscript": true, "object": true, "option": true, "optgroup": true,
	"output": true, "portal": true, "progress": true, "script": true,
	"select": true, "source": true, "style": true, "svg": true, "template": true,
	"textarea": true, "track": true, "video": true,
}

// allowedAttributes are inert cosmetic attributes. Anything not listed here is
// dropped, including all href/src variants and every event handler. Keep the
// list deliberately small — a screenshot does not need interactivity.
var allowedAttributes = map[string]bool{
	"align": true, "alt": true, "border": true, "class": true, "colspan": true,
	"datetime": true, "dir": true, "headers": true, "id": true, "lang": true,
	"rowspan": true, "scope": true, "start": true, "summary": true,
	"title": true, "valign": true, "width": true, "height": true,
	// `style` is handled separately through CSS().
}