package renderer

import "strings"

// Elements

func Html(props ...Renderer) *Element       { return El("html", props...) }
func Head(props ...Renderer) *Element       { return El("head", props...) }
func Link(props ...Renderer) *Element       { return El("link", props...) }
func Body(props ...Renderer) *Element       { return El("body", props...) }
func Title(props ...Renderer) *Element      { return El("title", props...) }
func Style(props ...Renderer) *Element      { return El("style", props...) }
func Header(props ...Renderer) *Element     { return El("header", props...) }
func Main(props ...Renderer) *Element       { return El("main", props...) }
func Article(props ...Renderer) *Element    { return El("article", props...) }
func P(props ...Renderer) *Element          { return El("p", props...) }
func Div(props ...Renderer) *Element        { return El("div", props...) }
func Blockquote(props ...Renderer) *Element { return El("blockquote", props...) }
func Ul(props ...Renderer) *Element         { return El("ul", props...) }
func Li(props ...Renderer) *Element         { return El("li", props...) }
func Dd(props ...Renderer) *Element         { return El("dd", props...) }
func A(props ...Renderer) *Element          { return El("a", props...) }
func Strong(props ...Renderer) *Element     { return El("strong", props...) }
func U(props ...Renderer) *Element          { return El("u", props...) }
func Ins(props ...Renderer) *Element        { return El("ins", props...) }
func Small(props ...Renderer) *Element      { return El("small", props...) }
func Sub(props ...Renderer) *Element        { return El("sub", props...) }
func Code(props ...Renderer) *Element       { return El("code", props...) }
func Samp(props ...Renderer) *Element       { return El("samp", props...) }
func Q(props ...Renderer) *Element          { return El("q", props...) }
func Rt(props ...Renderer) *Element         { return El("rt", props...) }
func Br(props ...Renderer) *Element         { return El("br", props...) }
func Bdi(props ...Renderer) *Element        { return El("bdi", props...) }
func Caption(props ...Renderer) *Element    { return El("caption", props...) }
func Td(props ...Renderer) *Element         { return El("td", props...) }
func Thead(props ...Renderer) *Element      { return El("thead", props...) }
func Tbody(props ...Renderer) *Element      { return El("tbody", props...) }
func Col(props ...Renderer) *Element        { return El("col", props...) }
func Picture(props ...Renderer) *Element    { return El("picture", props...) }
func Figcaption(props ...Renderer) *Element { return El("figcaption", props...) }
func Video(props ...Renderer) *Element      { return El("video", props...) }
func Track(props ...Renderer) *Element      { return El("track", props...) }
func Object(props ...Renderer) *Element     { return El("object", props...) }
func Iframe(props ...Renderer) *Element     { return El("iframe", props...) }
func Abbr(props ...Renderer) *Element       { return El("abbr", props...) }
func Meter(props ...Renderer) *Element      { return El("meter", props...) }
func Form(props ...Renderer) *Element       { return El("form", props...) }
func Input(props ...Renderer) *Element      { return El("input", props...) }
func Select(props ...Renderer) *Element     { return El("select", props...) }
func Option(props ...Renderer) *Element     { return El("option", props...) }
func Label(props ...Renderer) *Element      { return El("label", props...) }
func Datalist(props ...Renderer) *Element   { return El("datalist", props...) }
func Output(props ...Renderer) *Element     { return El("output", props...) }
func Command(props ...Renderer) *Element    { return El("command", props...) }
func Basefont(props ...Renderer) *Element   { return El("basefont", props...) }
func Center(props ...Renderer) *Element     { return El("center", props...) }
func Font(props ...Renderer) *Element       { return El("font", props...) }
func Frameset(props ...Renderer) *Element   { return El("frameset", props...) }
func Strike(props ...Renderer) *Element     { return El("strike", props...) }
func Btn(props ...Renderer) *Element        { return El("button", props...) }
func H1(props ...Renderer) *Element         { return El("h1", props...) }
func H2(props ...Renderer) *Element         { return El("h2", props...) }
func H3(props ...Renderer) *Element         { return El("h3", props...) }
func H4(props ...Renderer) *Element         { return El("h4", props...) }
func H5(props ...Renderer) *Element         { return El("h5", props...) }
func H6(props ...Renderer) *Element         { return El("h6", props...) }
func Pre(props ...Renderer) *Element        { return El("pre", props...) }
func Hr(props ...Renderer) *Element         { return El("hr", props...) }
func Ol(props ...Renderer) *Element         { return El("ol", props...) }
func Dt(props ...Renderer) *Element         { return El("dt", props...) }
func Span(props ...Renderer) *Element       { return El("span", props...) }
func Em(props ...Renderer) *Element         { return El("em", props...) }
func B(props ...Renderer) *Element          { return El("b", props...) }
func S(props ...Renderer) *Element          { return El("s", props...) }
func Mark(props ...Renderer) *Element       { return El("mark", props...) }
func Sup(props ...Renderer) *Element        { return El("sup", props...) }
func Dfn(props ...Renderer) *Element        { return El("dfn", props...) }
func Var(props ...Renderer) *Element        { return El("var", props...) }
func Kbd(props ...Renderer) *Element        { return El("kbd", props...) }
func Cite(props ...Renderer) *Element       { return El("cite", props...) }
func Ruby(props ...Renderer) *Element       { return El("ruby", props...) }
func Rp(props ...Renderer) *Element         { return El("rp", props...) }
func Wbr(props ...Renderer) *Element        { return El("wbr", props...) }
func Bdo(props ...Renderer) *Element        { return El("bdo", props...) }
func Table(props ...Renderer) *Element      { return El("table", props...) }
func Tr(props ...Renderer) *Element         { return El("tr", props...) }
func Th(props ...Renderer) *Element         { return El("th", props...) }
func Tfoot(props ...Renderer) *Element      { return El("tfoot", props...) }
func Colgroup(props ...Renderer) *Element   { return El("colgroup", props...) }
func Img(props ...Renderer) *Element        { return El("img", props...) }
func Figure(props ...Renderer) *Element     { return El("figure", props...) }
func MapElem(props ...Renderer) *Element    { return El("map", props...) }
func Area(props ...Renderer) *Element       { return El("area", props...) }
func Audio(props ...Renderer) *Element      { return El("audio", props...) }
func Source(props ...Renderer) *Element     { return El("source", props...) }
func Script(props ...Renderer) *Element     { return El("script", props...) }
func Noscript(props ...Renderer) *Element   { return El("noscript", props...) }
func Param(props ...Renderer) *Element      { return El("param", props...) }
func Embed(props ...Renderer) *Element      { return El("embed", props...) }
func Canvas(props ...Renderer) *Element     { return El("canvas", props...) }
func Address(props ...Renderer) *Element    { return El("address", props...) }
func Progress(props ...Renderer) *Element   { return El("progress", props...) }

// Attributes

func As(value string) *HtmlAttribute            { return Attr("as", value) }
func Accept(value string) *HtmlAttribute        { return Attr("accept", value) }
func AcceptCharset(value string) *HtmlAttribute { return Attr("accept-charset", value) }
func AccessKey(value string) *HtmlAttribute     { return Attr("accesskey", value) }
func Action(value string) *HtmlAttribute        { return Attr("action", value) }
func Align(value string) *HtmlAttribute         { return Attr("align", value) }
func Alt(value string) *HtmlAttribute           { return Attr("alt", value) }
func Aria(name, value string) *HtmlAttribute {
	if strings.HasPrefix(name, ":") {
		return Attr(":aria-"+name[1:], value)
	}
	return Attr("aria-"+name, value)
}
func Async(value string) *HtmlAttribute           { return Attr("async", value) }
func Autocomplete(value string) *HtmlAttribute    { return Attr("autocomplete", value) }
func Autofocus(value string) *HtmlAttribute       { return Attr("autofocus", value) }
func Autoplay(value string) *HtmlAttribute        { return Attr("autoplay", value) }
func BGColor(value string) *HtmlAttribute         { return Attr("bgcolor", value) }
func Border(value string) *HtmlAttribute          { return Attr("border", value) }
func Charset(value string) *HtmlAttribute         { return Attr("charset", value) }
func Checked(value string) *HtmlAttribute         { return Attr("checked", value) }
func CiteAttr(value string) *HtmlAttribute        { return Attr("cite", value) }
func Class(value string) *HtmlAttribute           { return Attr("class", value) }
func Color(value string) *HtmlAttribute           { return Attr("color", value) }
func Cols(value string) *HtmlAttribute            { return Attr("cols", value) }
func Colspan(value string) *HtmlAttribute         { return Attr("colspan", value) }
func Content(value string) *HtmlAttribute         { return Attr("content", value) }
func ContentEditable(value string) *HtmlAttribute { return Attr("contenteditable", value) }
func Controls(value string) *HtmlAttribute        { return Attr("controls", value) }
func Coords(value string) *HtmlAttribute          { return Attr("coords", value) }
func Data(name, value string) *HtmlAttribute {
	if strings.HasPrefix(name, ":") {
		return Attr(":data-"+name[1:], value)
	}
	return Attr("data-"+name, value)
}
func Datetime(value string) *HtmlAttribute     { return Attr("datetime", value) }
func Default(value string) *HtmlAttribute      { return Attr("default", value) }
func Defer(value string) *HtmlAttribute        { return Attr("defer", value) }
func Dir(value string) *HtmlAttribute          { return Attr("dir", value) }
func Dirname(value string) *HtmlAttribute      { return Attr("dirname", value) }
func Disabled(value string) *HtmlAttribute     { return Attr("disabled", value) }
func Download(value string) *HtmlAttribute     { return Attr("download", value) }
func Draggable(value string) *HtmlAttribute    { return Attr("draggable", value) }
func Enctype(value string) *HtmlAttribute      { return Attr("enctype", value) }
func EnterKeyHint(value string) *HtmlAttribute { return Attr("enterkeyhint", value) }
func For(value string) *HtmlAttribute          { return Attr("for", value) }
func FormAttr(value string) *HtmlAttribute     { return Attr("form", value) }
func FormAction(value string) *HtmlAttribute   { return Attr("formaction", value) }
func Headers(value string) *HtmlAttribute      { return Attr("headers", value) }
func Height(value string) *HtmlAttribute       { return Attr("height", value) }
func Hidden(value string) *HtmlAttribute       { return Attr("hidden", value) }
func High(value string) *HtmlAttribute         { return Attr("high", value) }
func Href(value string) *HtmlAttribute         { return Attr("href", value) }
func HrefLang(value string) *HtmlAttribute     { return Attr("hreflang", value) }
func HttpEquiv(value string) *HtmlAttribute    { return Attr("http-equiv", value) }
func ID(value string) *HtmlAttribute           { return Attr("id", value) }
func Inert(value string) *HtmlAttribute        { return Attr("inert", value) }
func InputMode(value string) *HtmlAttribute    { return Attr("inputmode", value) }
func IsMap(value string) *HtmlAttribute        { return Attr("ismap", value) }
func Kind(value string) *HtmlAttribute         { return Attr("kind", value) }
func LabelAttr(value string) *HtmlAttribute    { return Attr("label", value) }
func Src(value string) *HtmlAttribute          { return Attr("src", value) }
func Role(value string) *HtmlAttribute         { return Attr("role", value) }
func Lang(value string) *HtmlAttribute         { return Attr("lang", value) }
func List(value string) *HtmlAttribute         { return Attr("list", value) }
func Loop(value string) *HtmlAttribute         { return Attr("loop", value) }
func Low(value string) *HtmlAttribute          { return Attr("low", value) }
func Max(value string) *HtmlAttribute          { return Attr("max", value) }
func MaxLength(value string) *HtmlAttribute    { return Attr("maxlength", value) }
func Media(value string) *HtmlAttribute        { return Attr("media", value) }
func Method(value string) *HtmlAttribute       { return Attr("method", value) }
func Min(value string) *HtmlAttribute          { return Attr("min", value) }
func Multiple(value string) *HtmlAttribute     { return Attr("multiple", value) }
func Muted(value string) *HtmlAttribute        { return Attr("muted", value) }
func Name(value string) *HtmlAttribute         { return Attr("name", value) }
func NoValidate(value string) *HtmlAttribute   { return Attr("novalidate", value) }
func Type(value string) *HtmlAttribute         { return Attr("type", value) }
func Rel(value string) *HtmlAttribute          { return Attr("rel", value) }
func Width(value string) *HtmlAttribute        { return Attr("width", value) }
