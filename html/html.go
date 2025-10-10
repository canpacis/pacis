package html

// Elements

/*
# Doctype

In HTML, the doctype is the required <!doctype html> preamble found at the top of all documents. Its sole purpose is to prevent a browser from switching into so-called "quirks mode" when rendering a document; that is, the <!doctype html> doctype ensures that the browser makes a best-effort attempt at following the relevant specifications, rather than using a different rendering mode that is incompatible with some specifications.

The doctype is case-insensitive. The convention of MDN code examples is to use lowercase, but it's also common to write it as <!DOCTYPE html>.

https://developer.mozilla.org/en-US/docs/Glossary/Doctype
*/
var Doctype = VoidEl("!DOCTYPE", Attr("html", ""))

/*
# The Anchor element

The <a> HTML element (or anchor element), with its href attribute, creates a hyperlink to web pages, files, email addresses, locations in the same page, or anything else a URL can address.

Content within each <a> should indicate the link's destination. If the href attribute is present, pressing the enter key while focused on the <a> element will activate it.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a
*/
func A(items ...Item) *Element { return El("a", items...) }

/*
# The Abbreviation element

The <abbr> HTML element represents an abbreviation or acronym.

When including an abbreviation or acronym, provide a full expansion of the term in plain text on first use, along with the <abbr> to mark up the abbreviation. This informs the user what the abbreviation or acronym means.

The optional title attribute can provide an expansion for the abbreviation or acronym when a full expansion is not present. This provides a hint to user agents on how to announce/display the content while informing all users what the abbreviation means. If present, title must contain this full description and nothing else.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/abbr
*/
func Abbr(items ...Item) *Element { return El("abbr", items...) }

/*
# The Contact Address element

The <address> HTML element indicates that the enclosed HTML provides contact information for a person or people, or for an organization.

The contact information provided by an <address> element's contents can take whatever form is appropriate for the context, and may include any type of contact information that is needed, such as a physical address, URL, email address, phone number, social media handle, geographic coordinates, and so forth. The <address> element should include the name of the person, people, or organization to which the contact information refers.

<address> can be used in a variety of contexts, such as providing a business's contact information in the page header, or indicating the author of an article by including an <address> element within the <article>.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/address
*/
func Address(items ...Item) *Element { return El("address", items...) }

/*
# The Image Map Area element

The <area> HTML element defines an area inside an image map that has predefined clickable areas. An image map allows geometric areas on an image to be associated with hypertext links.

This element is used only within a <map> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/area
*/
func Area(items ...Item) *Element { return VoidEl("area", items...) }

/*
# The Article Contents element

The <article> HTML element represents a self-contained composition in a document, page, application, or site, which is intended to be independently distributable or reusable (e.g., in syndication). Examples include: a forum post, a magazine or newspaper article, or a blog entry, a product card, a user-submitted comment, an interactive widget or gadget, or any other independent item of content.

A given document can have multiple articles in it; for example, on a blog that shows the text of each article one after another as the reader scrolls, each post would be contained in an <article> element, possibly with one or more <section>s within.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/article
*/
func Article(items ...Item) *Element { return El("article", items...) }

/*
# The Aside element

The <aside> HTML element represents a portion of a document whose content is only indirectly related to the document's main content. Asides are frequently presented as sidebars or call-out boxes.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/aside
*/
func Aside(items ...Item) *Element { return El("aside", items...) }

/*
# The Embed Audio element

The <audio> HTML element is used to embed sound content in documents. It may contain one or more audio sources, represented using the src attribute or the <source> element: the browser will choose the most suitable one. It can also be the destination for streamed media, using a MediaStream (https://developer.mozilla.org/en-US/docs/Web/API/MediaStream).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/audio
*/
func Audio(items ...Item) *Element { return El("audio", items...) }

/*
# The Bring Attention To element

The <b> HTML element is used to draw the reader's attention to the element's contents, which are not otherwise granted special importance. This was formerly known as the Boldface element, and most browsers still draw the text in boldface. However, you should not use <b> for styling text or granting importance. If you wish to create boldface text, you should use the CSS font-weight property. If you wish to indicate an element is of special importance, you should use the <strong> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/b
*/
func B(items ...Item) *Element { return El("b", items...) }

/*
# The Document Base URL element

The <base> HTML element specifies the base URL to use for all relative URLs in a document. There can be only one <base> element in a document.

A document's used base URL can be accessed by scripts with Node.baseURI. If the document has no <base> elements, then baseURI defaults to location.href.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/base
*/
func Base(items ...Item) *Element { return VoidEl("base", items...) }

/*
# The Bidirectional Isolate element

The <bdi> HTML element tells the browser's bidirectional algorithm to treat the text it contains in isolation from its surrounding text. It's particularly useful when a website dynamically inserts some text and doesn't know the directionality of the text being inserted.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdi
*/
func Bdi(items ...Item) *Element { return El("bdi", items...) }

/*
# The Bidirectional Text Override element

The <bdo> HTML element overrides the current directionality of text, so that the text within is rendered in a different direction.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/bdo
*/
func Bdo(items ...Item) *Element { return El("bdo", items...) }

/*
# The Block Quotation element

The <blockquote> HTML element indicates that the enclosed text is an extended quotation. Usually, this is rendered visually by indentation (see Notes for how to change it). A URL for the source of the quotation may be given using the cite attribute, while a text representation of the source can be given using the <cite> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote
*/
func Blockquote(items ...Item) *Element { return El("blockquote", items...) }

/*
# The Document Body element

The <body> HTML element represents the content of an HTML document. There can be only one <body> element in a document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/body
*/
func Body(items ...Item) *Element { return El("body", items...) }

/*
# The Line Break element

The <br> HTML element produces a line break in text (carriage-return). It is useful for writing a poem or an address, where the division of lines is significant.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/br
*/
func Br(items ...Item) *Element { return VoidEl("br", items...) }

/*
# The Button element

The <button> HTML element is an interactive element activated by a user with a mouse, keyboard, finger, voice command, or other assistive technology. Once activated, it then performs an action, such as submitting a form or opening a dialog.

By default, HTML buttons are presented in a style resembling the platform the user agent runs on, but you can change buttons' appearance with CSS.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/button
*/
func Button(items ...Item) *Element { return El("button", items...) }

/*
# The Graphics Canvas element

Use the HTML <canvas> element with either the canvas scripting API or the WebGL API to draw graphics and animations.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/canvas
*/
func Canvas(items ...Item) *Element { return El("canvas", items...) }

/*
# The Table Caption element

The <caption> HTML element specifies the caption (or title) of a table, providing the table an accessible name or accessible description.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/caption
*/
func Caption(items ...Item) *Element { return El("caption", items...) }

/*
# The Citation element

The <cite> HTML element is used to mark up the title of a creative work. The reference may be in an abbreviated form according to context-appropriate conventions related to citation metadata.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/cite
*/
func Cite(items ...Item) *Element { return El("cite", items...) }

/*
# The Inline Code element

The <code> HTML element displays its contents styled in a fashion intended to indicate that the text is a short fragment of computer code. By default, the content text is displayed using the user agent's default monospace font.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/code
*/
func Code(items ...Item) *Element { return El("code", items...) }

/*
# The Table Column element

The <col> HTML element defines one or more columns in a column group represented by its parent <colgroup> element. The <col> element is only valid as a child of a <colgroup> element that has no span attribute defined.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/col
*/
func Col(items ...Item) *Element { return VoidEl("col", items...) }

/*
# The Table Column Group element

The <colgroup> HTML element defines a group of columns within a table.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/colgroup
*/
func Colgroup(items ...Item) *Element { return El("colgroup", items...) }

/*
# The Data element

The <data> HTML element links a given piece of content with a machine-readable translation. If the content is time- or date-related, the <time> element must be used.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/data
*/
func DataEl(items ...Item) *Element { return El("data", items...) }

/*
# The HTML Data List element

> Limited availability

The <datalist> HTML element contains a set of <option> elements that represent the permissible or recommended options available to choose from within other controls.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/datalist
*/
func Datalist(items ...Item) *Element { return El("datalist", items...) }

/*
# The Description Details element

The <dd> HTML element provides the description, definition, or value for the preceding term (<dt>) in a description list (<dl>).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dd
*/
func Dd(items ...Item) *Element { return El("dd", items...) }

/*
# The Deleted Text element

The <del> HTML element represents a range of text that has been deleted from a document. This can be used when rendering "track changes" or source code diff information, for example. The <ins> element can be used for the opposite purpose: to indicate text that has been added to the document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/del
*/
func Del(items ...Item) *Element { return El("del", items...) }

/*
# The Details disclosure element

The <details> HTML element creates a disclosure widget in which information is visible only when the widget is toggled into an open state. A summary or label must be provided using the <summary> element.

A disclosure widget is typically presented onscreen using a small triangle that rotates (or twists) to indicate open/closed state, with a label next to the triangle. The contents of the <summary> element are used as the label for the disclosure widget. The contents of the <details> provide the accessible description for the <summary>.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/details
*/
func Details(items ...Item) *Element { return El("details", items...) }

/*
# The Definition element

The <dfn> HTML element indicates a term to be defined. The <dfn> element should be used in a complete definition statement, where the full definition of the term can be one of the following:
  - The ancestor paragraph (a block of text, sometimes marked by a <p> element)
  - The <dt>/<dd> pairing
  - The nearest section ancestor of the <dfn> element

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dfn
*/
func Dfn(items ...Item) *Element { return El("dfn", items...) }

/*
# The Dialog element

The <dialog> HTML element represents a modal or non-modal dialog box or other interactive component, such as a dismissible alert, inspector, or subwindow.

The HTML <dialog> element is used to create both modal and non-modal dialog boxes. Modal dialog boxes interrupt interaction with the rest of the page being inert, while non-modal dialog boxes allow interaction with the rest of the page.

JavaScript should be used to display the <dialog> element. Use the .showModal() method to display a modal dialog and the .show() method to display a non-modal dialog. The dialog box can be closed using the .close() method or using the dialog method when submitting a <form> that is nested within the <dialog> element. Modal dialogs can also be closed by pressing the Esc key.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dialog
*/
func Dialog(items ...Item) *Element { return El("dialog", items...) }

/*
# The Content Division element

The <div> HTML element is the generic container for flow content. It has no effect on the content or layout until styled in some way using CSS (e.g., styling is directly applied to it, or some kind of layout model like Flexbox is applied to its parent element).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/div
*/
func Div(items ...Item) *Element { return El("div", items...) }

/*
# The Description List element

The <dl> HTML element represents a description list. The element encloses a list of groups of terms (specified using the <dt> element) and descriptions (provided by <dd> elements). Common uses for this element are to implement a glossary or to display metadata (a list of key-value pairs).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dl
*/
func Dl(items ...Item) *Element { return El("dl", items...) }

/*
# The Description Term element

The <dt> HTML element specifies a term in a description or definition list, and as such must be used inside a <dl> element. It is usually followed by a <dd> element; however, multiple <dt> elements in a row indicate several terms that are all defined by the immediate next <dd> element.

The subsequent <dd> (Description Details) element provides the definition or other related text associated with the term specified using <dt>.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/dt
*/
func Dt(items ...Item) *Element { return El("dt", items...) }

/*
# The Emphasis element

The <em> HTML element marks text that has stress emphasis. The <em> element can be nested, with each level of nesting indicating a greater degree of emphasis.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/em
*/
func Em(items ...Item) *Element { return El("em", items...) }

/*
# The Embed External Content element

The <embed> HTML element embeds external content at the specified point in the document. This content is provided by an external application or other source of interactive content such as a browser plug-in.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/embed
*/
func Embed(items ...Item) *Element { return VoidEl("embed", items...) }

/*
# The Field Set element

The <fieldset> HTML element is used to group several controls as well as labels (<label>) within a web form.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/fieldset
*/
func Fieldset(items ...Item) *Element { return El("fieldset", items...) }

/*
# The Figure Caption element

The <figcaption> HTML element represents a caption or legend describing the rest of the contents of its parent <figure> element, providing the <figure> an accessible name.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figcaption
*/
func Figcaption(items ...Item) *Element { return El("figcaption", items...) }

/*
# The Figure with Optional Caption element

The <figure> HTML element represents self-contained content, potentially with an optional caption, which is specified using the <figcaption> element. The figure, its caption, and its contents are referenced as a single unit.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/figure
*/
func Figure(items ...Item) *Element { return El("figure", items...) }

/*
# The Footer element

The <footer> HTML element represents a footer for its nearest ancestor sectioning content or sectioning root element. A <footer> typically contains information about the author of the section, copyright data or links to related documents.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/footer
*/
func Footer(items ...Item) *Element { return El("footer", items...) }

/*
# The Form element

The <form> HTML element represents a document section containing interactive controls for submitting information.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form
*/
func Form(items ...Item) *Element { return El("form", items...) }

/*
# The HTML Section Heading elements

The <h1> to <h6> HTML elements represent six levels of section headings. <h1> is the highest section level and <h6> is the lowest. By default, all heading elements create a block-level box in the layout, starting on a new line and taking up the full width available in their containing block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements
*/
func H1(items ...Item) *Element { return El("h1", items...) }

/*
# The HTML Section Heading elements

The <h1> to <h6> HTML elements represent six levels of section headings. <h1> is the highest section level and <h6> is the lowest. By default, all heading elements create a block-level box in the layout, starting on a new line and taking up the full width available in their containing block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements
*/
func H2(items ...Item) *Element { return El("h2", items...) }

/*
# The HTML Section Heading elements

The <h1> to <h6> HTML elements represent six levels of section headings. <h1> is the highest section level and <h6> is the lowest. By default, all heading elements create a block-level box in the layout, starting on a new line and taking up the full width available in their containing block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements
*/
func H3(items ...Item) *Element { return El("h3", items...) }

/*
# The HTML Section Heading elements

The <h1> to <h6> HTML elements represent six levels of section headings. <h1> is the highest section level and <h6> is the lowest. By default, all heading elements create a block-level box in the layout, starting on a new line and taking up the full width available in their containing block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements
*/
func H4(items ...Item) *Element { return El("h4", items...) }

/*
# The HTML Section Heading elements

The <h1> to <h6> HTML elements represent six levels of section headings. <h1> is the highest section level and <h6> is the lowest. By default, all heading elements create a block-level box in the layout, starting on a new line and taking up the full width available in their containing block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements
*/
func H5(items ...Item) *Element { return El("h5", items...) }

/*
# The HTML Section Heading elements

The <h1> to <h6> HTML elements represent six levels of section headings. <h1> is the highest section level and <h6> is the lowest. By default, all heading elements create a block-level box in the layout, starting on a new line and taking up the full width available in their containing block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements
*/
func H6(items ...Item) *Element { return El("h6", items...) }

/*
# The Document Metadata (Header) element

The <head> HTML element contains machine-readable information (metadata) about the document, like its title, scripts, and style sheets. There can be only one <head> element in an HTML document.

	> <head> primarily holds information for machine processing, not human-readability. For human-visible information, like top-level headings and listed authors, see the <header> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/head
*/
func Head(items ...Item) *Element { return El("head", items...) }

/*
# The Header element

The <header> HTML element represents introductory content, typically a group of introductory or navigational aids. It may contain some heading elements but also a logo, a search form, an author name, and other elements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/header
*/
func Header(items ...Item) *Element { return El("header", items...) }

/*
# The Heading Group element

The <hgroup> HTML element represents a heading and related content. It groups a single <h1>–<h6> element with one or more <p>.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hgroup
*/
func Hgroup(items ...Item) *Element { return El("hgroup", items...) }

/*
# The Thematic Break (Horizontal Rule) element

The <hr> HTML element represents a thematic break between paragraph-level elements: for example, a change of scene in a story, or a shift of topic within a section.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/hr
*/
func Hr(items ...Item) *Element { return VoidEl("hr", items...) }

/*
# The HTML Document / Root element

The <html> HTML element represents the root (top-level element) of an HTML document, so it is also referred to as the root element. All other elements must be descendants of this element. There can be only one <html> element in a document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/html
*/
func Html(items ...Item) *Element { return El("html", items...) }

/*
# The Idiomatic Text element

The <i> HTML element represents a range of text that is set off from the normal text for some reason, such as idiomatic text, technical terms, taxonomical designations, among others. Historically, these have been presented using italicized type, which is the original source of the <i> naming of this element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/i
*/
func I(items ...Item) *Element { return El("i", items...) }

/*
# The Inline Frame element

The <iframe> HTML element represents a nested browsing context, embedding another HTML page into the current one.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/iframe
*/
func Iframe(items ...Item) *Element { return El("iframe", items...) }

/*
# The Image Embed element

The <img> HTML element embeds an image into the document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img
*/
func Img(items ...Item) *Element { return VoidEl("img", items...) }

/*
# The HTML Input element

The <input> HTML element is used to create interactive controls for web-based forms in order to accept data from the user; a wide variety of types of input data and control widgets are available, depending on the device and user agent. The <input> element is one of the most powerful and complex in all of HTML due to the sheer number of combinations of input types and attributes.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input
*/
func Input(items ...Item) *Element { return VoidEl("input", items...) }

/*
# The Inserted Text element

The <ins> HTML element represents a range of text that has been added to a document. You can use the <del> element to similarly represent a range of text that has been deleted from the document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ins
*/
func Ins(items ...Item) *Element { return El("ins", items...) }

/*
# The Keyboard Input element

The <kbd> HTML element represents a span of inline text denoting textual user input from a keyboard, voice input, or any other text entry device. By convention, the user agent defaults to rendering the contents of a <kbd> element using its default monospace font, although this is not mandated by the HTML standard.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/kbd
*/
func Kbd(items ...Item) *Element { return El("kbd", items...) }

/*
# The Label element

The <label> HTML element represents a caption for an item in a user interface.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/label
*/
func Label(items ...Item) *Element { return El("label", items...) }

/*
# The Field Set Legend element

The <legend> HTML element represents a caption for the content of its parent <fieldset>.
In customizable <select> elements, the <legend> element is allowed as a child of <optgroup>, to provide a label that is easy to target and style. This replaces any text set in the <optgroup> element's label attribute, and it has the same semantics.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/legend
*/
func Legend(items ...Item) *Element { return El("legend", items...) }

/*
# The List Item element

The <li> HTML element is used to represent an item in a list. It must be contained in a parent element: an ordered list (<ol>), an unordered list (<ul>), or a menu (<menu>). In menus and unordered lists, list items are usually displayed using bullet points. In ordered lists, they are usually displayed with an ascending counter on the left, such as a number or letter.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/li
*/
func Li(items ...Item) *Element { return El("li", items...) }

/*
# The External Resource Link element

The <link> HTML element specifies relationships between the current document and an external resource.
This element is most commonly used to link to stylesheets, but is also used to establish site icons (both "favicon" style icons and icons for the home screen and apps on mobile devices) among other things.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/link
*/
func Link(items ...Item) *Element { return VoidEl("link", items...) }

/*
# The Main element

The <main> HTML element represents the dominant content of the <body> of a document. The main content area consists of content that is directly related to or expands upon the central topic of a document, or the central functionality of an application.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/main
*/
func Main(items ...Item) *Element { return El("main", items...) }

/*
# The Image Map element

The <map> HTML element is used with <area> elements to define an image map (a clickable link area).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/map
*/
func MapEl(items ...Item) *Element { return El("map", items...) }

/*
# The Mark Text element

The <mark> HTML element represents text which is marked or highlighted for reference or notation purposes due to the marked passage's relevance in the enclosing context.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/mark
*/
func Mark(items ...Item) *Element { return El("mark", items...) }

/*
# The Menu element

The <menu> HTML element is described in the HTML specification as a semantic alternative to <ul>, but treated by browsers (and exposed through the accessibility tree) as no different than <ul>. It represents an unordered list of items (which are represented by <li> elements).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/menu
*/
func Menu(items ...Item) *Element { return El("menu", items...) }

/*
# The metadata element

The <meta> HTML element represents metadata that cannot be represented by other meta-related elements, such as <base>, <link>, <script>, <style>, or <title>.
The type of metadata provided by the <meta> element can be one of the following:

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta
*/
func Meta(items ...Item) *Element { return VoidEl("meta", items...) }

/*
# The HTML Meter element

The <meter> HTML element represents either a scalar value within a known range or a fractional value.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meter
*/
func Meter(items ...Item) *Element { return El("meter", items...) }

/*
# The Navigation Section element

The <nav> HTML element represents a section of a page whose purpose is to provide navigation links, either within the current document or to other documents. Common examples of navigation sections are menus, tables of contents, and indexes.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/nav
*/
func Nav(items ...Item) *Element { return El("nav", items...) }

/*
# The Noscript element

The <noscript> HTML element defines a section of HTML to be inserted if a script type on the page is unsupported or if scripting is currently turned off in the browser.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/noscript
*/
func Noscript(items ...Item) *Element { return El("noscript", items...) }

/*
# The External Object element

The <object> HTML element represents an external resource, which can be treated as an image, a nested browsing context, or a resource to be handled by a plugin.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/object
*/
func Object(items ...Item) *Element { return El("object", items...) }

/*
# The Ordered List element

The <ol> HTML element represents an ordered list of items — typically rendered as a numbered list.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ol
*/
func Ol(items ...Item) *Element { return El("ol", items...) }

/*
# The Option Group element

The <optgroup> HTML element creates a grouping of options within a <select> element.
In customizable <select> elements, the <legend> element is allowed as a child of <optgroup>, to provide a label that is easy to target and style. This replaces any text set in the <optgroup> element's label attribute, and it has the same semantics.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/optgroup
*/
func Optgroup(items ...Item) *Element { return El("optgroup", items...) }

/*
# The HTML Option element

The <option> HTML element is used to define an item contained in a <select>, an <optgroup>, or a <datalist> element. As such, <option> can represent menu items in popups and other lists of items in an HTML document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/option
*/
func Option(items ...Item) *Element { return El("option", items...) }

/*
# The Output element

The <output> HTML element is a container element into which a site or app can inject the results of a calculation or the outcome of a user action.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/output
*/
func Output(items ...Item) *Element { return El("output", items...) }

/*
# The Paragraph element

The <p> HTML element represents a paragraph. Paragraphs are usually represented in visual media as blocks of text separated from adjacent blocks by blank lines and/or first-line indentation, but HTML paragraphs can be any structural grouping of related content, such as images or form fields.
Paragraphs are block-level elements, and notably will automatically close if another block-level element is parsed before the closing </p> tag. See "Tag omission" below.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/p
*/
func P(items ...Item) *Element { return El("p", items...) }

/*
# The Picture element

The <picture> HTML element contains zero or more <source> elements and one <img> element to offer alternative versions of an image for different display/device scenarios.
The browser will consider each child <source> element and choose the best match among them. If no matches are found—or the browser doesn't support the <picture> element—the URL of the <img> element's src attribute is selected. The selected image is then presented in the space occupied by the <img> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/picture
*/
func Picture(items ...Item) *Element { return El("picture", items...) }

/*
# The Preformatted Text element

The <pre> HTML element represents preformatted text which is to be presented exactly as written in the HTML file. The text is typically rendered using a non-proportional, or monospaced font.

Whitespace inside this element is displayed as written, with one exception. If one or more leading newline characters are included immediately following the opening <pre> tag, the first newline character is stripped.

<pre> elements' text content is parsed as HTML, so if you want to ensure that your text content stays as plain text, some syntax characters, such as <, may need to be escaped using their respective character references. See escaping ambiguous characters for more information.

<pre> elements commonly contain <code>, <samp>, and <kbd> elements, to represent computer code, computer output, and user input, respectively.

By default, <pre> is a block-level element, i.e., its default display value is block.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/pre
*/
func Pre(items ...Item) *Element { return El("pre", items...) }

/*
# The Progress Indicator element

The <progress> HTML element displays an indicator showing the completion progress of a task, typically displayed as a progress bar.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/progress
*/
func Progress(items ...Item) *Element { return El("progress", items...) }

/*
# The Inline Quotation element

The <q> HTML element indicates that the enclosed text is a short inline quotation. Most modern browsers implement this by surrounding the text in quotation marks. This element is intended for short quotations that don't require paragraph breaks; for long quotations use the <blockquote> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/q
*/
func Q(items ...Item) *Element { return El("q", items...) }

/*
# The Ruby Fallback Parenthesis element

The <rp> HTML element is used to provide fall-back parentheses for browsers that do not support display of ruby annotations using the <ruby> element. One <rp> element should enclose each of the opening and closing parentheses that wrap the <rt> element that contains the annotation's text.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rp
*/
func Rp(items ...Item) *Element { return El("rp", items...) }

/*
# The Ruby Text element

The <rt> HTML element specifies the ruby text component of a ruby annotation, which is used to provide pronunciation, translation, or transliteration information for East Asian typography. The <rt> element must always be contained within a <ruby> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/rt
*/
func Rt(items ...Item) *Element { return El("rt", items...) }

/*
# The Ruby Annotation element

The <ruby> HTML element represents small annotations that are rendered above, below, or next to base text, usually used for showing the pronunciation of East Asian characters. It can also be used for annotating other kinds of text, but this usage is less common.

The term ruby originated as a unit of measurement used by typesetters, representing the smallest size that text can be printed on newsprint while remaining legible.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ruby
*/
func Ruby(items ...Item) *Element { return El("ruby", items...) }

/*
# The Strikethrough element

The <s> HTML element renders text with a strikethrough, or a line through it. Use the <s> element to represent things that are no longer relevant or no longer accurate. However, <s> is not appropriate when indicating document edits; for that, use the <del> and <ins> elements, as appropriate.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/s
*/
func S(items ...Item) *Element { return El("s", items...) }

/*
# The Sample Output element

The <samp> HTML element is used to enclose inline text which represents sample (or quoted) output from a computer program. Its contents are typically rendered using the browser's default monospaced font (such as Courier or Lucida Console).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/samp
*/
func Samp(items ...Item) *Element { return El("samp", items...) }

/*
# The Script element

The <script> HTML element is used to embed executable code or data; this is typically used to embed or refer to JavaScript code. The <script> element can also be used with other languages, such as WebGL's GLSL shader programming language and JSON.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/script
*/
func Script(items ...Item) *Element { return El("script", items...) }

/*
# The generic search element

The <search> HTML element is a container representing the parts of the document or application with form controls or other content related to performing a search or filtering operation. The <search> element semantically identifies the purpose of the element's contents as having search or filtering capabilities. The search or filtering functionality can be for the website or application, the current web page or document, or the entire Internet or subsection thereof.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/search
*/
func Search(items ...Item) *Element { return El("search", items...) }

/*
# The Generic Section element

The <section> HTML element represents a generic standalone section of a document, which doesn't have a more specific semantic element to represent it. Sections should always have a heading, with very few exceptions.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/section
*/
func Section(items ...Item) *Element { return El("section", items...) }

/*
# The HTML Select element

The <select> HTML element represents a control that provides a menu of options.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/select
*/
func Select(items ...Item) *Element { return El("select", items...) }

/*
# The Web Component Slot element

The <slot> HTML element—part of the Web Components technology suite—is a placeholder inside a web component that you can fill with your own markup, which lets you create separate DOM trees and present them together.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/slot
*/
func Slot(items ...Item) *Element { return El("slot", items...) }

/*
# The Side Comment element

The <small> HTML element represents side-comments and small print, like copyright and legal text, independent of its styled presentation. By default, it renders text within it one font-size smaller, such as from small to x-small.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/small
*/
func Small(items ...Item) *Element { return El("small", items...) }

/*
# The Media or Image Source element

The <source> HTML element specifies one or more media resources for the <picture>, <audio>, and <video> elements. It is a void element, which means that it has no content and does not require a closing tag. This element is commonly used to offer the same media content in multiple file formats in order to provide compatibility with a broad range of browsers given their differing support for image file formats and media file formats.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/source
*/
func Source(items ...Item) *Element { return VoidEl("source", items...) }

/*
# The Content Span element

The <span> HTML element is a generic inline container for phrasing content, which does not inherently represent anything. It can be used to group elements for styling purposes (using the class or id attributes), or because they share attribute values, such as lang. It should be used only when no other semantic element is appropriate. <span> is very much like a <div> element, but <div> is a block-level element whereas a <span> is an inline-level element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/span
*/
func Span(items ...Item) *Element { return El("span", items...) }

/*
# The Strong Importance element

The <strong> HTML element indicates that its contents have strong importance, seriousness, or urgency. Browsers typically render the contents in bold type.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/strong
*/
func Strong(items ...Item) *Element { return El("strong", items...) }

/*
# The Style Information element

The <style> HTML element contains style information for a document, or part of a document. It contains CSS, which is applied to the contents of the document containing the <style> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/style
*/
func Style(items ...Item) *Element { return El("style", items...) }

/*
# The Subscript element

The <sub> HTML element specifies inline text which should be displayed as subscript for solely typographical reasons. Subscripts are typically rendered with a lowered baseline using smaller text.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sub
*/
func Sub(items ...Item) *Element { return El("sub", items...) }

/*
# The Disclosure Summary element

The <summary> HTML element specifies a summary, caption, or legend for a <details> element's disclosure box. Clicking the <summary> element toggles the state of the parent <details> element open and closed.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/summary
*/
func Summary(items ...Item) *Element { return El("summary", items...) }

/*
# The Superscript element

The <sup> HTML element specifies inline text which is to be displayed as superscript for solely typographical reasons. Superscripts are usually rendered with a raised baseline using smaller text.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/sup
*/
func Sup(items ...Item) *Element { return El("sup", items...) }

/*
# The Table element

The <table> HTML element represents tabular data—that is, information presented in a two-dimensional table comprised of rows and columns of cells containing data.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/table
*/
func Table(items ...Item) *Element { return El("table", items...) }

/*
# The Table Body element

The <tbody> HTML element encapsulates a set of table rows (<tr> elements), indicating that they comprise the body of a table's (main) data.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tbody
*/
func Tbody(items ...Item) *Element { return El("tbody", items...) }

/*
# The Table Data Cell element

The <td> HTML element defines a cell of a table that contains data and may be used as a child of the <tr> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/td
*/
func Td(items ...Item) *Element { return El("td", items...) }

/*
# The Content Template element

The <template> HTML element serves as a mechanism for holding HTML fragments, which can either be used later via JavaScript or generated immediately into shadow DOM.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template
*/
func Template(items ...Item) *Element { return El("template", items...) }

/*
# The Textarea element

The <textarea> HTML element represents a multi-line plain-text editing control, useful when you want to allow users to enter a sizeable amount of free-form text, for example a comment on a review or feedback form.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/textarea
*/
func Textarea(items ...Item) *Element { return El("textarea", items...) }

/*
# The Table Foot element

The <tfoot> HTML element encapsulates a set of table rows (<tr> elements), indicating that they comprise the foot of a table with information about the table's columns. This is usually a summary of the columns, e.g., a sum of the given numbers in a column.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tfoot
*/
func Tfoot(items ...Item) *Element { return El("tfoot", items...) }

/*
# The Table Header element

The <th> HTML element defines a cell as the header of a group of table cells and may be used as a child of the <tr> element. The exact nature of this group is defined by the scope and headers attributes.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/th
*/
func Th(items ...Item) *Element { return El("th", items...) }

/*
# The Table Head element

The <thead> HTML element encapsulates a set of table rows (<tr> elements), indicating that they comprise the head of a table with information about the table's columns. This is usually in the form of column headers (<th> elements).

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/thead
*/
func Thead(items ...Item) *Element { return El("thead", items...) }

/*
# The (Date) Time element

The <time> HTML element represents a specific period in time. It may include the datetime attribute to translate dates into machine-readable format, allowing for better search engine results or custom features such as reminders.

It may represent one of the following:
  - A time on a 24-hour clock.
  - A precise date in the Gregorian calendar (with optional time and timezone information).
  - A valid time duration.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/time
*/
func Time(items ...Item) *Element { return El("time", items...) }

/*
# The Document Title element

The <title> HTML element defines the document's title that is shown in a browser's title bar or a page's tab. It only contains text; HTML tags within the element, if any, are also treated as plain text.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/title
*/
func Title(items ...Item) *Element { return El("title", items...) }

/*
# The Table Row element

The <tr> HTML element defines a row of cells in a table. The row's cells can then be established using a mix of <td> (data cell) and <th> (header cell) elements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/tr
*/
func Tr(items ...Item) *Element { return El("tr", items...) }

/*
# The Embed Text Track element

The <track> HTML element is used as a child of the media elements, <audio> and <video>.
Each track element lets you specify a timed text track (or time-based data) that can be displayed in parallel with the media element, for example to overlay subtitles or closed captions on top of a video or alongside audio tracks.

Multiple tracks can be specified for a media element, containing different kinds of timed text data, or timed text data that has been translated for different locales.
The data that is used will either be the track that has been set to be the default, or a kind and translation based on user preferences.

The tracks are formatted in WebVTT format (.vtt files) — Web Video Text Tracks.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/track
*/
func Track(items ...Item) *Element { return VoidEl("track", items...) }

/*
# The Unarticulated Annotation (Underline) element

The <u> HTML element represents a span of inline text which should be rendered in a way that indicates that it has a non-textual annotation. This is rendered by default as a single solid underline, but may be altered using CSS.

> This element used to be called the "Underline" element in older versions of HTML, and is still sometimes misused in this way. To underline text, you should instead apply a style that includes the CSS text-decoration property set to underline.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/u
*/
func U(items ...Item) *Element { return El("u", items...) }

/*
# The Unordered List element

The <ul> HTML element represents an unordered list of items, typically rendered as a bulleted list.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/ul
*/
func Ul(items ...Item) *Element { return El("ul", items...) }

/*
# The Variable element

The <var> HTML element represents the name of a variable in a mathematical expression or a programming context. It's typically presented using an italicized version of the current typeface, although that behavior is browser-dependent.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/var
*/
func Var(items ...Item) *Element { return El("var", items...) }

/*
# The Video Embed element

The <video> HTML element embeds a media player which supports video playback into the document. You can use <video> for audio content as well, but the <audio> element may provide a more appropriate user experience.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/video
*/
func Video(items ...Item) *Element { return El("video", items...) }

/*
# The Line Break Opportunity element

The <wbr> HTML element represents a word break opportunity—a position within text where the browser may optionally break a line, though its line-breaking rules would not otherwise create a break at that location.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/wbr
*/
func Wbr(items ...Item) *Element { return VoidEl("wbr", items...) }

// Attributes

func As(value string) *Attribute { return Attr("as", value) }

/*
# HTML attribute: accept

The accept attribute takes as its value a comma-separated list of one or more file types, or unique file type specifiers, describing which file types to allow.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/accept
*/
func Accept(value string) *Attribute        { return Attr("accept", value) }
func AcceptCharset(value string) *Attribute { return Attr("accept-charset", value) }
func AccessKey(value string) *Attribute     { return Attr("accesskey", value) }
func Action(value string) *Attribute        { return Attr("action", value) }
func Align(value string) *Attribute         { return Attr("align", value) }

/*
# <img> alt attribute

Defines text that can replace the image in the page.

> Browsers do not always display images. There are a number of situations in which a browser might not display images, such as:

  - Non-visual browsers (such as those used by people with visual impairments)
  - The user chooses not to display images (saving bandwidth, privacy reasons)
  - The image is invalid or an unsupported type

In these cases, the browser may replace the image with the text in the element's alt attribute. For these reasons and others, provide a useful value for alt whenever possible.

Setting this attribute to an empty string (alt="") indicates that this image is not a key part of the content (it's decoration or a tracking pixel), and that non-visual browsers may omit it from rendering. Visual browsers will also hide the broken image icon if the alt attribute is empty and the image failed to display.

This attribute is also used when copying and pasting the image to text, or saving a linked image to a bookmark.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/img#alt
*/
func Alt(value string) *Attribute        { return Attr("alt", value) }
func Aria(name, value string) *Attribute { return Attr("aria-"+name, value) }
func Async(value string) *Attribute      { return Attr("async", value) }

/*
# HTML attribute: autocomplete

It is available on <input> elements that take a text or numeric value as input, <textarea> elements, <select> elements, and <form> elements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/autocomplete
*/
func Autocomplete(value string) *Attribute { return Attr("autocomplete", value) }

/*
# HTML attribute: autofocus
*/
func Autofocus(value string) *Attribute { return Attr("autofocus", value) }

/*
# HTML attribute: autoplay
*/
func Autoplay(value string) *Attribute { return Attr("autoplay", value) }

/*
# HTML attribute: charset

This attribute declares the document's character encoding. If the attribute is present, its value must be an ASCII case-insensitive match for the string "utf-8", because UTF-8 is the only valid encoding for HTML5 documents. <meta> elements which declare a character encoding must be located entirely within the first 1024 bytes of the document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta#charset
*/
func Charset(value string) *Attribute { return Attr("charset", value) }

var Checked = Attr("checked", "")

/*
# <blockquote> cite attribute

A URL that designates a source document or message for the information quoted. This attribute is intended to point to information explaining the context or the reference for the quote.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/blockquote#cite
*/
func CiteAttr(value string) *Attribute { return Attr("cite", value) }

/*
# HTML class global attribute

The class global attribute is a list of the classes of the element, separated by ASCII whitespace.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/class
*/
func Class(class string) *Attribute     { return Attr("class", class) }
func ColorAttr(value string) *Attribute { return Attr("color", value) }
func Cols(value string) *Attribute      { return Attr("cols", value) }
func Colspan(value string) *Attribute   { return Attr("colspan", value) }

/*
# HTML attribute: content

The content attribute specifies the value of a metadata name defined by the <meta> name attribute.
It takes a string as its value, and the expected syntax varies depending on the name value used.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/content
*/
func Content(value string) *Attribute { return Attr("content", value) }

/*
# HTML contenteditable global attribute

The contenteditable global attribute is an enumerated attribute indicating if the element should be editable by the user. If so, the browser modifies its widget to allow editing.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/contenteditable
*/
var ContentEditable = Attr("contenteditable", "")

func Controls(value string) *Attribute { return Attr("controls", value) }
func Coords(value string) *Attribute   { return Attr("coords", value) }

/*
# HTML attribute: crossorigin

The crossorigin attribute, valid on the <audio>, <img>, <link>, <script>, and <video> elements, provides support for CORS, defining how the element handles cross-origin requests, thereby enabling the configuration of the CORS requests for the element's fetched data. Depending on the element, the attribute can be a CORS settings attribute.

The crossorigin content attribute on media elements is a CORS settings attribute.

These attributes are enumerated, and have the following possible values:

Request uses CORS headers and credentials flag is set to 'same-origin'. There is no exchange of user credentials via cookies, client-side TLS certificates or HTTP authentication, unless destination is the same origin.

Request uses CORS headers, credentials flag is set to 'include' and user credentials are always included.

Setting the attribute name to an empty value, like crossorigin or crossorigin="", is the same as anonymous.

An invalid keyword and an empty string will be handled as the anonymous keyword.

By default (that is, when the attribute is not specified), CORS is not used at all. The user agent will not ask for permission for full access to the resource and in the case of a cross-origin request, certain limitations will be applied based on the type of element concerned:

	> The crossorigin attribute is not supported for rel="icon" in Chromium-based browsers. See the open Chromium issue.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/crossorigin
*/
func Crossorigin(value string) *Attribute { return Attr("crossorigin", value) }

/*
# HTML data-* global attribute

The data-* global attributes form a class of attributes called custom data attributes, that allow proprietary information to be exchanged between the HTML and its DOM representation by scripts.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/data-*
*/
func Data(name, value string) *Attribute  { return Attr("data-"+name, value) }
func Datetime(value string) *Attribute    { return Attr("datetime", value) }
func DefaultAttr(value string) *Attribute { return Attr("default", value) }

var Defer = Attr("defer", "")

/*
# HTML attribute: dirname

The dirname attribute can be used on the <textarea> element and several <input> types and describes the directionality of the element's text content during form submission. The browser uses this attribute's value to determine whether text the user has entered is left-to-right or right-to-left oriented. When used, the element's text directionality value is included in form submission data along with the dirname attribute's value as the name of the field.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/dirname
*/
func Dirname(value string) *Attribute { return Attr("dirname", value) }

/*
# HTML attribute: disabled

The Boolean disabled attribute, when present, makes the element not mutable, focusable, or even submitted with the form. The user can neither edit nor focus on the control, nor its form control descendants.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/disabled
*/
var Disabled = Attr("disabled", "")

func Download(value string) *Attribute { return Attr("download", value) }

/*
# HTML draggable global attribute

The draggable global attribute is an enumerated attribute that indicates whether the element can be dragged, either with native browser behavior or the HTML Drag and Drop API.

The draggable attribute may be applied to elements that strictly fall under the HTML namespace, which means that it cannot be applied to SVGs.
For more information about what namespace declarations look like, and what they do, see Namespace crash course.

draggable can have the following values:

	> This attribute is enumerated and not Boolean. A value of true or false is mandatory, and shorthand like <img draggable> is forbidden. The correct usage is <img draggable="true">.

If this attribute is not set, its default value is auto, which means drag behavior is the default browser behavior: only text selections, images, and links can be dragged. For other elements, the event ondragstart must be set for drag and drop to work, as shown in this comprehensive example.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/draggable
*/
func Draggable(value string) *Attribute    { return Attr("draggable", value) }
func Enctype(value string) *Attribute      { return Attr("enctype", value) }
func EnterKeyHint(value string) *Attribute { return Attr("enterkeyhint", value) }

/*
# HTML attribute: for

The for attribute is an allowed attribute for <label> and <output>. When used on a <label> element it indicates the form element that this label describes. When used on an <output> element it allows for an explicit relationship between the elements that represent values which are used in the output.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/for
*/
func For(value string) *Attribute { return Attr("for", value) }

/*
# HTML attribute: form

The form HTML attribute associates a form-associated element with a <form> element within the same document. This attribute applies to the <button>, <fieldset>, <input>, <object>, <output>, <select>, and <textarea> elements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/form
*/
func FormAttr(value string) *Attribute { return Attr("form", value) }

/*
# <form> action attribute

The URL that processes the form submission. This value can be overridden by a formaction attribute on a <button>, <input type="submit">, or <input type="image"> element. This attribute is ignored when method="dialog" is set.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form#action
*/
func FormAction(value string) *Attribute { return Attr("formaction", value) }
func Headers(value string) *Attribute    { return Attr("headers", value) }

/*
# SVG attribute: height

The height attribute defines the vertical length of an element in the user coordinate system.

https://developer.mozilla.org/en-US/docs/Web/SVG/Reference/Attribute/height
*/
func Height(value string) *Attribute { return Attr("height", value) }

/*
# HTML hidden global attribute

The hidden global attribute is an enumerated attribute indicating that the browser should not render the contents of the element. For example, it can be used to hide elements of the page that can't be used until the login process has been completed.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/hidden
*/
func Hidden(value string) *Attribute { return Attr("hidden", value) }
func High(value string) *Attribute   { return Attr("high", value) }

/*
# <a> href attribute

The URL that the hyperlink points to. Links are not restricted to HTTP-based URLs — they can use any URL scheme supported by browsers:

  - Telephone numbers with tel: URLs

  - Email addresses with mailto: URLs

  - SMS text messages with sms: URLs

  - Executable code with javascript: URLs

  - While web browsers may not support other URL schemes, websites can with registerProtocolHandler()

Moreover other URL features can locate specific parts of the resource, including:

  - Sections of a page with document fragments

  - Specific text portions with text fragments

  - Pieces of media files with media fragments

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a#href
*/
func Href(value string) *Attribute { return Attr("href", value) }

/*
# <a> hreflang attribute

Hints at the human language of the linked URL. No built-in functionality. Allowed values are the same as the global lang attribute.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a#hreflang
*/
func HrefLang(value string) *Attribute { return Attr("hreflang", value) }

/*
# <meta> http-equiv attribute

The http-equiv attribute of the <meta> element allows you to provide processing instructions for the browser as if the response that returned the document included certain HTTP headers.
The metadata is document-level metadata that applies to the whole page.

When a <meta> element has an http-equiv attribute, a content attribute defines the corresponding http-equiv value.
For example, the following <meta> tag tells the browser to refresh the page after 5 minutes:

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta/http-equiv
*/
func HttpEquiv(value string) *Attribute { return Attr("http-equiv", value) }

/*
# HTML id global attribute

The id global attribute defines an identifier (ID) that must be unique within the entire document.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/id
*/
func ID(value string) *Attribute        { return Attr("id", value) }
func Inert(value string) *Attribute     { return Attr("inert", value) }
func InputMode(value string) *Attribute { return Attr("inputmode", value) }
func IsMap(value string) *Attribute     { return Attr("ismap", value) }
func Kind(value string) *Attribute      { return Attr("kind", value) }
func LabelAttr(value string) *Attribute { return Attr("label", value) }
func Src(value string) *Attribute       { return Attr("src", value) }
func Role(value string) *Attribute      { return Attr("role", value) }

/*
# HTML lang global attribute

The lang global attribute helps define the language of an element: the language that non-editable elements are written in, or the language that the editable elements should be written in by the user. The attribute contains a single BCP 47 language tag.

	> The default value of lang is the empty string, which means that the language is unknown. Therefore, it is recommended to always specify an appropriate value for this attribute.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/lang
*/
func Lang(value string) *Attribute { return Attr("lang", value) }
func List(value string) *Attribute { return Attr("list", value) }
func Loop(value string) *Attribute { return Attr("loop", value) }
func Low(value string) *Attribute  { return Attr("low", value) }

/*
# HTML attribute: max

The max attribute defines the maximum value that is acceptable and valid for the input containing the attribute. If the value of the element is greater than this, the element fails validation. This value must be greater than or equal to the value of the min attribute. If the max attribute is present but is not specified or is invalid, no max value is applied. If the max attribute is valid and a non-empty value is greater than the maximum allowed by the max attribute, constraint validation will prevent form submission.

The max attribute is valid for the numeric input types, including the date, month, week, time, datetime-local, number and range types, and both the <progress> and <meter> elements. It is a number that specifies the most positive value a form control to be considered valid.

If the value exceeds the max value allowed, the validityState.rangeOverflow will be true, and the control will be matched by the :out-of-range and :invalid pseudo-classes.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/max
*/
func Max(value string) *Attribute { return Attr("max", value) }

/*
# HTML attribute: maxlength

The maxlength attribute defines the maximum string length that the user can enter into an <input> or <textarea>. The attribute must have an integer value of 0 or higher.

The length is measured in UTF-16 code units, which is often but not always equal to the number of characters. If no maxlength is specified, or an invalid value is specified, the input has no maximum length.

Any maxlength value must be greater than or equal to the value of minlength, if present and valid. The input will fail constraint validation if the length of the text value of the field is greater than maxlength UTF-16 code units long. Constraint validation is only applied when the value is changed by the user.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/maxlength
*/
func MaxLength(value string) *Attribute { return Attr("maxlength", value) }

/*
# HTML attribute: media

The media attribute defines which media the theme color defined in the content attribute should be applied to. Its value is a media query, which defaults to all if the attribute is missing. This attribute is only relevant when the element's name attribute is set to theme-color. Otherwise, it has no effect, and should not be included.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta#media
*/
func Media(value string) *Attribute { return Attr("media", value) }

/*
# HTML attribute: method

The HTTP method to submit the form with. The only allowed methods/values are (case insensitive):

post: The POST method; form data sent as the request body.
get (default): The GET; form data appended to the action URL with a ? separator. Use this method when the form has no side effects.
dialog: When the form is inside a <dialog>, closes the dialog and causes a submit event to be fired on submission, without submitting data or clearing the form.
This value is overridden by formmethod attributes on <button>, <input type="submit">, or <input type="image"> elements.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form#method
*/
func Method(value string) *Attribute { return Attr("method", value) }

/*
# HTML attribute: min

The min attribute defines the minimum value that is acceptable and valid for the input containing the attribute. If the value of the element is less than this, the element fails validation. This value must be less than or equal to the value of the max attribute.

Some input types have a default minimum. If the input has no default minimum and a value is specified for min that can't be converted to a valid number (or no minimum value is set), the input has no minimum value.

It is valid for the input types including: date, month, week, time, datetime-local, number and range types, and the <meter> element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/min
*/
func Min(value string) *Attribute { return Attr("min", value) }

/*
# HTML attribute: minlength

The minlength attribute defines the minimum string length that the user can enter into an <input> or <textarea>. The attribute must have an integer value of 0 or higher.

The length is measured in UTF-16 code units, which is often but not always equal to the number of characters. If no minlength is specified, or an invalid value is specified, the input has no minimum length. This value must be less than or equal to the value of maxlength, otherwise the value will never be valid, as it is impossible to meet both criteria.

The input will fail constraint validation if the length of the text value of the field is less than minlength UTF-16 code units long, with validityState.tooShort returning true. Constraint validation is only applied when the value is changed by the user. Once submission fails, some browsers will display an error message indicating the minimum length required and the current length.

minlength does not imply required; an input only violates a minlength constraint if the user has input a value. If an input is not required, an empty string can be submitted even if minlength is set.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/minlength
*/
func MinLength(value string) *Attribute { return Attr("minlength", value) }

/*
# HTML attribute: multiple

The Boolean multiple attribute, if set, means the form control accepts one or more values. The attribute is valid for the email and file input types and the <select>. The manner by which the user opts for multiple values depends on the form control.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/multiple
*/
func Multiple(value string) *Attribute { return Attr("multiple", value) }

/*
# HTML attribute: muted
*/
func Muted(value string) *Attribute { return Attr("muted", value) }

/*
# HTML attribute: name

The name and content attributes can be used together to provide document metadata in terms of name-value pairs, with the name attribute giving the metadata name, and the content attribute giving the value.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/meta#name
*/
func Name(value string) *Attribute { return Attr("name", value) }

/*
# HTML attribute: novalidate

This Boolean attribute indicates that the form shouldn't be validated when submitted. If this attribute is not set (and therefore the form is validated), it can be overridden by a formnovalidate attribute on a <button>, <input type="submit">, or <input type="image"> element belonging to the form.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/form#novalidate
*/
var NoValidate = Attr("novalidate", "")

/*
# HTML attribute: placeholder

Effective placeholder text includes a word or short phrase that hints at the expected data type, not an explanation or prompt. The placeholder must not be used instead of a <label>. As the placeholder is not visible if the value of the form control is not null, using placeholder instead of a <label> for a prompt harms usability and accessibility.

The placeholder attribute is supported by the following input types: text, search, url, tel, email, and password. It is also supported by the <textarea> element. The example below shows the placeholder attribute in use to explain the expected format of an input field.

	> Except in <textarea> elements, the placeholder attribute can't include any line feeds (LF) or carriage returns (CR). If either is included in the value, the placeholder text will be clipped.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/placeholder
*/
func Placeholder(value string) *Attribute { return Attr("placeholder", value) }

/*
# HTML style global attribute

The style global attribute contains CSS styling declarations to be applied to the element. Note that it is recommended for styles to be defined in a separate file or files. This attribute and the <style> element have mainly the purpose of allowing for quick styling, for example for testing purposes.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/style
*/
func StyleAttr(value string) *Attribute { return Attr("style", value) }

/*
# HTML tabindex global attribute

The tabindex global attribute allows developers to make HTML elements focusable, allow or prevent them from being sequentially focusable (usually with the Tab key, hence the name) and determine their relative ordering for sequential focus navigation.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/tabindex
*/
func Tabindex(value string) *Attribute { return Attr("tabindex", value) }

/*
# HTML title global attribute

The title global attribute contains text representing advisory information related to the element it belongs to.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/title
*/
func TitleAttr(value string) *Attribute { return Attr("title", value) }

/*
# <input> type attribute

A string specifying the type of control to render. For example, to create a checkbox, a value of checkbox is used. If omitted (or an unknown value is specified), the input type text is used, creating a plaintext input field.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input#type
*/
func Type(value string) *Attribute { return Attr("type", value) }

/*
# HTML attribute: rel

The rel attribute defines the relationship between a linked resource and the current document. Valid on <link>, <a>, <area>, and <form>, the supported values depend on the element on which the attribute is found.

The type of relationships is given by the value of the rel attribute, which, if present, must have a value that is an unordered set of unique space-separated keywords. Differently from a class name, which does not express semantics, the rel attribute must express tokens that are semantically valid for both machines and humans. The current registries for the possible values of the rel attribute are the IANA link relation registry, the HTML Living Standard, and the freely-editable existing-rel-values page in the microformats wiki, as suggested by the Living Standard. If a rel attribute not present in one of the three sources above is used some HTML validators (such as the W3C Markup Validation Service) will generate a warning.

The following table lists some of the most important existing keywords. Every keyword within a space-separated value should be unique within that value.

The rel attribute is relevant to the <link>, <a>, <area>, and <form> elements, but some values only relevant to a subset of those elements. Like all HTML keyword attribute values, these values are case-insensitive.

The rel attribute has no default value. If the attribute is omitted or if none of the values in the attribute are supported, then the document has no particular relationship with the destination resource other than there being a hyperlink between the two. In this case, on <link> and <form>, if the rel attribute is absent, has no keywords, or if not one or more of the space-separated keywords above, then the element does not create any links. <a> and <area> will still created links, but without a defined relationship.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/rel
*/
func Rel(value string) *Attribute { return Attr("rel", value) }

/*
# HTML attribute: required

The Boolean required attribute, if present, indicates that the user must specify a value for the input before the owning form can be submitted.

The required attribute is supported by text, search, url, tel, email, password, date, month, week, time, datetime-local, number, checkbox, radio, file, <input> types along with the <select> and <textarea> form control elements. If present on any of these input types and elements, the :required pseudo class will match. If the attribute is not included, the :optional pseudo class will match.

The attribute is not supported on, or relevant to, range and color input types, as both have default values. Type color defaults to #000000. Type range defaults to the midpoint between min and max — with min and max defaulting to 0 and 100 respectively in most browsers if not declared. required is also not supported on the hidden input type — users cannot be expected to fill out a hidden form field. Finally, required is not supported on any button input types, including image.

In the case of a same named group of radio buttons, if a single radio button in the group has the required attribute, a radio button in that group must be checked, although it doesn't have to be the one on which the attribute is applied. To improve code maintenance, it is recommended to either include the required attribute in every same-named radio button in the group, or else in none.

In the case of a same named group of checkbox input types, only the checkboxes with the required attribute are required.

	> Setting aria-required="true" tells a screen reader that an element (any element) is required, but has no bearing on the optionality of the element.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Attributes/required
*/
var Required = Attr("rel", "")

/*
# SVG attribute: width

The width attribute defines the horizontal length of an element in the user coordinate system.

https://developer.mozilla.org/en-US/docs/Web/SVG/Reference/Attribute/width
*/
func Width(value string) *Attribute { return Attr("width", value) }

/*
# <input> value attribute

The input control's value. When specified in the HTML, this is the initial value, and from then on it can be altered or retrieved at any time using JavaScript to access the respective HTMLInputElement object's value property. The value attribute is always optional, though should be considered mandatory for checkbox, radio, and hidden.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/input#value
*/
func Value(value string) *Attribute { return Attr("value", value) }

/*
# <template> shadowrootmode attribute

Creates a shadow root for the parent element. It is a declarative version of the Element.attachShadow() method and accepts the same enumerated values.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/template#shadowrootmode
*/
func ShadowRootMode(value string) *Attribute { return Attr("shadowrootmode", value) }

/*
# HTML slot global attribute

The slot global attribute assigns a slot in a shadow DOM shadow tree to an element: An element with a slot attribute is assigned to the slot created by the <slot> element whose name attribute's value matches that slot attribute's value. You can have multiple elements assigned to the same slot by using the same slot name. Elements without a slot attribute are assigned to the unnamed slot, if one exists.

For examples, see our Using templates and slots guide.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Global_attributes/slot
*/
func SlotAttr(value string) *Attribute { return Attr("slot", value) }

/*
# <a> target attribute

Where to display the linked URL, as the name for a browsing context (a tab, window, or <iframe>). The following keywords have special meanings for where to load the URL:

  - _self: The current browsing context. (Default)
  - _blank: Usually a new tab, but users can configure browsers to open a new window instead.
  - _parent: The parent browsing context of the current one. If no parent, behaves as _self.
  - _top: The topmost browsing context. To be specific, this means the "highest" context that's an ancestor of the current one. If no ancestors, behaves as _self.
  - _unfencedTop: Allows embedded fenced frames to navigate the top-level frame (i.e., traversing beyond the root of the fenced frame, unlike other reserved destinations). Note that the navigation will still succeed if this is used outside of a fenced frame context, but it will not act like a reserved keyword.

https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/a#target
*/
func Target(value string) *Attribute       { return Attr("target", value) }
func PropertyAttr(value string) *Attribute { return Attr("property", value) }
