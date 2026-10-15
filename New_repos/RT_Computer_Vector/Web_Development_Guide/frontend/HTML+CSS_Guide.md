# HTML AND CSS GUIDE 

PARAGRAPH AND FUNDAMENTAL STRUCTURE 

**The first example**
```html
 <!DOCTYPE html> --> must decleration 
<html>
<head>
<title>Page Title</title>
</head>
<body>

<h1>This is a Heading</h1>
<p>This is a paragraph.</p>

</body>
</html> 
```

This is how HTML code looks; this is just here as an example before we start 

HTML follows a tag format kind of system, for example 
```html
 <tagname> Content goes here... </tagname> 
```

**HTML Heading tags** (for adding HTML headings)
```html
<h1>Page Title</h1>
<h2>Section</h2>
<h3>Subsection</h3>
<body>body content here</body>
<p>paragraph here</p>
<html>entire html between this</html>
<pre>this is pre formated text</pre>
<div>generic container</div>
```
Remember, these tags have entire wording between them cause they are the main formatters 

**Different Attributes for heading tag**
```html
<!-- changing the default font size --!>
<h1 style="font-size:60px;">Heading 1</h1>

<!-- unique identifier --!>
<h2 id="pricing">Pricing</h2>
<a href="#pricing">Jump to Pricing</a>

<!-- tool tip to hover over --!>
<h2 title="Click to learn more about pricing">Pricing</h2>

<!-- langauge for the heading text --!.
<h1 lang="fr">Bienvenue</h1>

<!-- text direction --!>
<h2 dir="rtl">مرحبا بكم</h2>

<!-- custom data attributes for JS --!>
<h2 data-section="pricing" data-analytics="true">Pricing</h2>

<!-- hides the heading --!>
<h2 hidden>This heading is hidden</h2>

<!-- make it focusable via keyboard --!>
<h2 tabindex="0">Focusable Heading</h2>

<!-- making content editable --!>
<h1 contenteditable="true">Click to edit me</h1>

<!-- making it accessible for screen readers --!>
<h2 aria-label="Pricing plans overview">Pricing</h2>
<div role="heading" aria-level="2">Custom element acting as an h2</div>
```
**Different Attributes for paragraph tag** 
```html
<!-- adding colour to a paragraph --!>
<p style="color:red;">This is a red paragraph.</p>  

<!-- title attributes !-->
 <p title="I'm a tooltip">This is a paragraph.</p>

<!-- a longer qouted paragrpgh form another source --!>
<blockquote cite="https://example.com/source">
  <p>This is a quoted paragraph from another source.</p>
</blockquote>

<!-- address tag specifically for the contact info --!>
<blockquote cite="https://example.com/source">
  <p>This is a quoted paragraph from another source.</p>
</blockquote>
```

**Link-related info**
```html
<!-- this is how to add a image --!>
<a href="https://www.w3schools.com">This is a link</a> 

<!-- this id for internal relative link --!>
<a href="/about.html">About Us</a> 

<!-- this is to jump in the same page --!>
<h2 id="contact">Contact</h2> 
<a href="#contact">Go to Contact section</a>

<!-- this is to open in a new tab --!>
<a href="https://example.com" target="_blank" rel="noopener noreferrer">
  Opens in new tab
</a>

<!-- link with a title --!>
<a href="https://example.com" title="just a example">Visit our HTML Tutorial</a>

<!-- button ad a link --!>
<button onclick="document.location='default.asp'">HTML Tutorial</button>

<!-- this is for email links --!>
<a href="mailto:someone@example.com">Email us</a>

<!-- this is for phone links --!>
<a href="tel:+911234567890">Call us</a>

<!-- Download instead of opening it --!>
<a href="tel:+911234567890">Call us</a>

<!-- link tag for inside head use --!>
<link rel="stylesheet" href="styles.css">
<link rel="icon" href="favicon.ico">

<!-- nav tag for semantic wrapper for group of links --!>
<nav>
  <a href="/">Home</a>
  <a href="/about.html">About</a>
  <a href="/contact.html">Contact</a>
</nav>
```
**HTML image**
```html
<!-- basic image --!>
<img src="photo.jpg" alt="Description of the image">

<!-- with adjusted size --!>
<img src="photo.jpg" alt="A sunset over the hills" width="600" height="400">

<!-- responsive image (dynamic size) --!>
<img
  src="photo-800w.jpg"
  srcset="photo-400w.jpg 400w, photo-800w.jpg 800w, photo-1200w.jpg 1200w"
  sizes="(max-width: 600px) 400px, 800px"
  alt="Responsive photo example">

<!-- serve different formats/art directions --!>
<picture>
  <source srcset="photo.avif" type="image/avif">
  <source srcset="photo.webp" type="image/webp">
  <img src="photo.jpg" alt="Fallback for browsers that don't support avif/webp">
</picture>

<!-- lazy loading images --!>
<img src="photo.jpg" alt="Lazy loaded image" loading="lazy">

<!-- image as links --!>
<a href="https://example.com">
  <img src="banner.jpg" alt="Click to visit Example">
</a>

<!-- image with a caption --!>
<figure>
  <img src="chart.jpg" alt="Sales growth chart">
  <figcaption>Fig. 1 — Sales growth, 2024–2025</figcaption>
</figure>

<!-- clickable regions on one image --!>
<img src="map.jpg" usemap="#worldmap" alt="World map">
<map name="worldmap">
  <area shape="rect" coords="0,0,100,100" href="asia.html" alt="Asia">
  <area shape="circle" coords="200,150,50" href="europe.html" alt="Europe">
</map>

<!-- favicon --!>
<link rel="icon" href="favicon.ico" type="image/x-icon">

<!-- image/logos for crispyness --!>
<svg width="100" height="100">
  <circle cx="50" cy="50" r="40" fill="blue" />
</svg>
```

**empty html tags**
tags with no ending tag; some, like links and images, are already mentioned above 

```html
<br> <!-- line break tag --!>
<base> <!-- setting a base url --!>
<hr> <!-- horizontal break --!>
<meta> <!-- page meta data --!>
<source> <!-- media source for <picture>, <audio>, <video> --!>
<track>   <!-- subtitle/caption track for <audio>/<video> -->
<wbr>     <!-- optional line-break opportunity in text -->
<input>   <!-- form input field -->
<col>     <!-- defines a column within a <colgroup> -->
<embed>   <!-- embeds external content (flash, plugins, etc.) -->
<area>    <!-- clickable region inside an image map -->
<base>    <!-- sets base URL for all relative links on the page -->
```
**Making bookmarks in HTML**
```html
method to create : 
<h2 id="C4">Chapter 4</h2>

method to jump:
 <a href="#C4">Jump to Chapter 4</a>

jumping on a bookmark on another page:
<a href="html_demo.html#C4">Jump to Chapter 4</a>
```

STYLING AND FORMATTING STUFF 

**Text formatting**
```html
<strong>bold text here</strong>
<mark>highlighted</mark>
<em>emphaisesed</em>
<q>done well</q>
<i>itallic</i>
<small>small text</small>
<del>deleted text</del>
<ins>insterted text</ins>
<sub>subscription text</sub>
<sup>superscript text</sup>
```
**HTML style attributes**
```html
<!-- basic syntax --!>
<p style="color: red; font-size: 18px;">Styled text</p>

<!-- text styling --!> 
<p style="color: navy; font-size: 20px; font-weight: bold; font-style: italic; text-align: center; text-decoration: underline; line-height: 1.5; letter-spacing: 1px;">
  Styled paragraph
</p>

<!-- background & color --!>
<div style="background-color: lightgray; color: black;">Box with background</div>
<div style="background-image: url('bg.jpg'); background-size: cover;">Image background</div>

<!-- spacing (box model) --!>
<div style="margin: 20px; padding: 10px; width: 300px; height: 150px;">
  Spaced box
</div>

<!-- Border and round corners --!>
<div style="border: 2px solid #333; border-radius: 10px; padding: 15px;">
  Bordered box
</div>

<!-- display and positioning --!>
<div style="display: flex; justify-content: center; align-items: center;">Flex box</div>
<div style="position: absolute; top: 10px; left: 20px;">Positioned box</div>

<!-- Shadow and opacity --!>
<div style="box-shadow: 2px 2px 8px rgba(0,0,0,0.3); opacity: 0.9;">
  Box with shadow
</div>
```
*Here are few more details*

Common colour adding tricks 
```html
<!-- 1. Named color (147 keywords built into CSS) -->
<p style="color: tomato;">Named color</p>

<!-- 2. Hex code -->
<p style="color: #ff6347;">Hex color</p>

<!-- 3. RGB -->
<p style="color: rgb(255, 99, 71);">RGB color</p>

<!-- 4. RGBA (RGB + transparency) -->
<p style="color: rgba(255, 99, 71, 0.5);">RGBA with 50% opacity</p>

<!-- 5. HSL (Hue, Saturation, Lightness) -->
<p style="color: hsl(9, 100%, 64%);">HSL color</p>

<!-- 6. HSLA (HSL + transparency) -->
<p style="color: hsla(9, 100%, 64%, 0.5);">HSLA with transparency</p>
```
handy trick with colors (currentColour)
```html
<div style="color: teal; border: 2px solid currentColor;">
  Border automatically matches the text color
</div>
```
**Different methods of adding CSS**
```html
<!-- 1. Inline CSS — written directly on the element via the style attribute -->
<p style="color: red; font-size: 18px;">Inline styled text</p>

<!-- 2. Internal CSS — written inside a <style> tag in the <head> -->
<head>
  <style>
    p {
      color: red;
      font-size: 18px;
    }
  </style>
</head>

<!-- 3. External CSS — written in a separate .css file, linked via <link> -->
<head>
  <link rel="stylesheet" href="styles.css">
</head>
```
Inline always wins over internal, and internal wins over external — regardless of which one appears last in the code. (You can still override this with !important, but that's generally a last resort, not a habit.)

> Methods 2 and 3 are superior when making or working with real project and big sites 

**Adding color to links**
will give an example to help you; this would make links far better and more decorative 
```html
<style>
a:link {
  color: green;
  background-color: transparent;
  text-decoration: none;
}

a:visited {
  color: pink;
  background-color: transparent;
  text-decoration: none;
}

a:hover {
  color: red;
  background-color: transparent;
  text-decoration: underline;
}

a:active {
  color: yellow;
  background-color: transparent;
  text-decoration: underline;
}
</style> 
```
> we keep using "a" cause "a" is used to define links 

here is one more example 
```html
<style>
a:link, a:visited {
  background-color: #f44336;
  color: white;
  padding: 15px 25px;
  text-align: center;
  text-decoration: none;
  display: inline-block;
}

a:hover, a:active {
  background-color: red;
}
</style> 
```
**HTML Quotation and citation element**
```html
<!-- abbreviation -->
<p>The <abbr title="World Health Organization">WHO</abbr> was founded in 1948.</p>

<!-- address -->
<address>
  Written by <a href="mailto:jane@example.com">Jane Doe</a><br>
  Visit us at Example St, Ahmedabad, India
</address>

<!-- bi-directional override -->
<bdo dir="rtl">This text will be displayed right to left</bdo>

<!-- blockquote (long, block-level quotation) -->
<blockquote cite="https://example.com/source">
  <p>This is a longer quoted passage taken from another source, displayed as its own indented block.</p>
</blockquote>

<!-- cite (title of a creative work) -->
<p>The Great Gatsby was written by F. Scott Fitzgerald.</p>

<!-- q (short inline quotation) -->
<p>She said, <q>this changes everything</q>, and walked out of the room.</p>
```
