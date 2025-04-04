package html

import "strings"

// Elements

func Html(props ...I) Element       { return El("html", props...) }
func Head(props ...I) Element       { return El("head", props...) }
func Link(props ...I) Element       { return El("link", props...) }
func Body(props ...I) Element       { return El("body", props...) }
func Title(props ...I) Element      { return El("title", props...) }
func Style(props ...I) Element      { return El("style", props...) }
func Header(props ...I) Element     { return El("header", props...) }
func Main(props ...I) Element       { return El("main", props...) }
func Article(props ...I) Element    { return El("article", props...) }
func P(props ...I) Element          { return El("p", props...) }
func Div(props ...I) Element        { return El("div", props...) }
func Blockquote(props ...I) Element { return El("blockquote", props...) }
func Ul(props ...I) Element         { return El("ul", props...) }
func Li(props ...I) Element         { return El("li", props...) }
func Dd(props ...I) Element         { return El("dd", props...) }
func A(props ...I) Element          { return El("a", props...) }
func Strong(props ...I) Element     { return El("strong", props...) }
func U(props ...I) Element          { return El("u", props...) }
func Ins(props ...I) Element        { return El("ins", props...) }
func Small(props ...I) Element      { return El("small", props...) }
func Sub(props ...I) Element        { return El("sub", props...) }
func Cde(props ...I) Element        { return El("code", props...) }
func Samp(props ...I) Element       { return El("samp", props...) }
func Q(props ...I) Element          { return El("q", props...) }
func Rt(props ...I) Element         { return El("rt", props...) }
func Br(props ...I) Element         { return El("br", props...) }
func Bdi(props ...I) Element        { return El("bdi", props...) }
func Caption(props ...I) Element    { return El("caption", props...) }
func Td(props ...I) Element         { return El("td", props...) }
func Thead(props ...I) Element      { return El("thead", props...) }
func Tbody(props ...I) Element      { return El("tbody", props...) }
func Col(props ...I) Element        { return El("col", props...) }
func Picture(props ...I) Element    { return El("picture", props...) }
func Figcaption(props ...I) Element { return El("figcaption", props...) }
func Video(props ...I) Element      { return El("video", props...) }
func Track(props ...I) Element      { return El("track", props...) }
func Object(props ...I) Element     { return El("object", props...) }
func Iframe(props ...I) Element     { return El("iframe", props...) }
func Abbr(props ...I) Element       { return El("abbr", props...) }
func Meter(props ...I) Element      { return El("meter", props...) }
func Form(props ...I) Element       { return El("form", props...) }
func Inpt(props ...I) Element       { return El("input", props...) }
func Select(props ...I) Element     { return El("select", props...) }
func Slot(props ...I) Element       { return El("slot", props...) }
func Option(props ...I) Element     { return El("option", props...) }
func Lbl(props ...I) Element        { return El("label", props...) }
func Datalist(props ...I) Element   { return El("datalist", props...) }
func Output(props ...I) Element     { return El("output", props...) }
func Command(props ...I) Element    { return El("command", props...) }
func Basefont(props ...I) Element   { return El("basefont", props...) }
func Center(props ...I) Element     { return El("center", props...) }
func Font(props ...I) Element       { return El("font", props...) }
func Frameset(props ...I) Element   { return El("frameset", props...) }
func Strike(props ...I) Element     { return El("strike", props...) }
func Btn(props ...I) Element        { return El("button", props...) }
func H1(props ...I) Element         { return El("h1", props...) }
func H2(props ...I) Element         { return El("h2", props...) }
func H3(props ...I) Element         { return El("h3", props...) }
func H4(props ...I) Element         { return El("h4", props...) }
func H5(props ...I) Element         { return El("h5", props...) }
func H6(props ...I) Element         { return El("h6", props...) }
func Pre(props ...I) Element        { return El("pre", props...) }
func Hr(props ...I) Element         { return El("hr", props...) }
func Ol(props ...I) Element         { return El("ol", props...) }
func Dt(props ...I) Element         { return El("dt", props...) }
func Span(props ...I) Element       { return El("span", props...) }
func Em(props ...I) Element         { return El("em", props...) }
func B(props ...I) Element          { return El("b", props...) }
func S(props ...I) Element          { return El("s", props...) }
func Mark(props ...I) Element       { return El("mark", props...) }
func Sup(props ...I) Element        { return El("sup", props...) }
func Dfn(props ...I) Element        { return El("dfn", props...) }
func Var(props ...I) Element        { return El("var", props...) }
func Kbd(props ...I) Element        { return El("kbd", props...) }
func Cite(props ...I) Element       { return El("cite", props...) }
func Ruby(props ...I) Element       { return El("ruby", props...) }
func Rp(props ...I) Element         { return El("rp", props...) }
func Wbr(props ...I) Element        { return El("wbr", props...) }
func Bdo(props ...I) Element        { return El("bdo", props...) }
func Table(props ...I) Element      { return El("table", props...) }
func Tr(props ...I) Element         { return El("tr", props...) }
func Th(props ...I) Element         { return El("th", props...) }
func Tfoot(props ...I) Element      { return El("tfoot", props...) }
func Template(props ...I) Element   { return El("template", props...) }
func Colgroup(props ...I) Element   { return El("colgroup", props...) }
func Img(props ...I) Element        { return El("img", props...) }
func Figure(props ...I) Element     { return El("figure", props...) }
func MapElem(props ...I) Element    { return El("map", props...) }
func Area(props ...I) Element       { return El("area", props...) }
func Audio(props ...I) Element      { return El("audio", props...) }
func Source(props ...I) Element     { return El("source", props...) }
func Script(props ...I) Element     { return El("script", props...) }
func Noscript(props ...I) Element   { return El("noscript", props...) }
func Param(props ...I) Element      { return El("param", props...) }
func Embed(props ...I) Element      { return El("embed", props...) }
func Canvas(props ...I) Element     { return El("canvas", props...) }
func Address(props ...I) Element    { return El("address", props...) }
func Progress(props ...I) Element   { return El("progress", props...) }
func Section(props ...I) Element    { return El("section", props...) }
func Aside(props ...I) Element      { return El("aside", props...) }
func Meta(props ...I) Element       { return El("meta", props...) }
func Footer(props ...I) Element     { return El("footer", props...) }

// Attributes

func As(value string) Attribute            { return Attr("as", value) }
func Accept(value string) Attribute        { return Attr("accept", value) }
func AcceptCharset(value string) Attribute { return Attr("accept-charset", value) }
func AccessKey(value string) Attribute     { return Attr("accesskey", value) }
func Action(value string) Attribute        { return Attr("action", value) }
func Align(value string) Attribute         { return Attr("align", value) }
func Alt(value string) Attribute           { return Attr("alt", value) }
func Aria(name, value string) Attribute {
	if strings.HasPrefix(name, ":") {
		return Attr(":aria-"+name[1:], value)
	}
	return Attr("aria-"+name, value)
}
func Async(value string) Attribute           { return Attr("async", value) }
func Autocomplete(value string) Attribute    { return Attr("autocomplete", value) }
func Autofocus(value string) Attribute       { return Attr("autofocus", value) }
func Autoplay(value string) Attribute        { return Attr("autoplay", value) }
func BGColor(value string) Attribute         { return Attr("bgcolor", value) }
func Border(value string) Attribute          { return Attr("border", value) }
func Charset(value string) Attribute         { return Attr("charset", value) }
func Checked(value string) Attribute         { return Attr("checked", value) }
func CiteAttr_(value string) Attribute       { return Attr("cite", value) }
func Class(value string) Attribute           { return Attr("class", value) }
func ColorAttr(value string) Attribute       { return Attr("color", value) }
func Cols(value string) Attribute            { return Attr("cols", value) }
func Colspan(value string) Attribute         { return Attr("colspan", value) }
func Content(value string) Attribute         { return Attr("content", value) }
func ContentEditable(value string) Attribute { return Attr("contenteditable", value) }
func Controls(value string) Attribute        { return Attr("controls", value) }
func Coords(value string) Attribute          { return Attr("coords", value) }
func Data(name, value string) Attribute {
	if strings.HasPrefix(name, ":") {
		return Attr(":data-"+name[1:], value)
	}
	return Attr("data-"+name, value)
}
func Datetime(value string) Attribute { return Attr("datetime", value) }
func Default(value string) Attribute  { return Attr("default", value) }

var Defer = Attr("defer")

func Dir(value string) Attribute            { return Attr("dir", value) }
func Dirname(value string) Attribute        { return Attr("dirname", value) }
func Disabled(value string) Attribute       { return Attr("disabled", value) }
func Download(value string) Attribute       { return Attr("download", value) }
func Draggable(value string) Attribute      { return Attr("draggable", value) }
func Enctype(value string) Attribute        { return Attr("enctype", value) }
func EnterKeyHint(value string) Attribute   { return Attr("enterkeyhint", value) }
func For(value string) Attribute            { return Attr("for", value) }
func FormAttr(value string) Attribute       { return Attr("form", value) }
func FormAction(value string) Attribute     { return Attr("formaction", value) }
func Headers(value string) Attribute        { return Attr("headers", value) }
func Height(value string) Attribute         { return Attr("height", value) }
func Hidden(value string) Attribute         { return Attr("hidden", value) }
func High(value string) Attribute           { return Attr("high", value) }
func Href(value string) Attribute           { return Attr("href", value) }
func HrefLang(value string) Attribute       { return Attr("hreflang", value) }
func HttpEquiv(value string) Attribute      { return Attr("http-equiv", value) }
func HtmlFor(value string) Attribute        { return Attr("htmlFor", value) }
func ID(value string) Attribute             { return Attr("id", value) }
func Inert(value string) Attribute          { return Attr("inert", value) }
func InputMode(value string) Attribute      { return Attr("inputmode", value) }
func IsMap(value string) Attribute          { return Attr("ismap", value) }
func Kind(value string) Attribute           { return Attr("kind", value) }
func LabelAttr(value string) Attribute      { return Attr("label", value) }
func Src(value string) Attribute            { return Attr("src", value) }
func Role(value string) Attribute           { return Attr("role", value) }
func Lang(value string) Attribute           { return Attr("lang", value) }
func List(value string) Attribute           { return Attr("list", value) }
func Loop(value string) Attribute           { return Attr("loop", value) }
func Low(value string) Attribute            { return Attr("low", value) }
func Max(value string) Attribute            { return Attr("max", value) }
func MaxLength(value string) Attribute      { return Attr("maxlength", value) }
func Media(value string) Attribute          { return Attr("media", value) }
func Method(value string) Attribute         { return Attr("method", value) }
func Min(value string) Attribute            { return Attr("min", value) }
func Multiple(value string) Attribute       { return Attr("multiple", value) }
func Muted(value string) Attribute          { return Attr("muted", value) }
func Name(value string) Attribute           { return Attr("name", value) }
func NoValidate(value string) Attribute     { return Attr("novalidate", value) }
func Placeholder(value string) Attribute    { return Attr("placeholder", value) }
func Type(value string) Attribute           { return Attr("type", value) }
func Rel(value string) Attribute            { return Attr("rel", value) }
func Width(value string) Attribute          { return Attr("width", value) }
func Value(value string) Attribute          { return Attr("value", value) }
func Shadowrootmode(value string) Attribute { return Attr("shadowrootmode", value) }
func SlotAttr(value string) Attribute       { return Attr("slot", value) }
func Target(value string) Attribute         { return Attr("target", value) }
