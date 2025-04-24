package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	c "github.com/canpacis/pacis/ui/components"
	"github.com/canpacis/pacis/ui/html"
)

func readattr(attr html.Attribute) string {
	var buf bytes.Buffer
	attr.Render(context.Background(), &buf)
	return buf.String()
}

type Link struct {
	html.Element
}

func (l *Link) Render(ctx context.Context, w io.Writer) error {
	hrefattr, ok := l.Element.GetAttribute("href")
	if !ok {
		return l.Element.Render(ctx, w)
	}
	href := readattr(hrefattr)

	// In-page link, just render regular anchor
	if strings.HasPrefix(href, "#") {
		return l.Element.Render(ctx, w)
	}

	parsed, err := url.Parse(href)
	// Valid external url or broken url, just render regular anchor
	if err != nil || len(parsed.Host) != 0 {
		l.Element.AddAttribute(html.Target("blank"))
		l.Element.AddAttribute(html.Rel("noreferer"))
		return l.Element.Render(ctx, w)
	}

	l.Element.AddAttribute(c.On("mouseenter", fmt.Sprintf("$prefetch.queue('%s')", href)))
	l.Element.AddAttribute(c.On("click", fmt.Sprintf("$prefetch.load('%s', $event)", href)))

	_, eager := l.Element.GetAttribute("eager")
	if eager {
		l.Element.RemoveAttribute("eager")
		l.Element.AddAttribute(c.X("intersect", fmt.Sprintf("$prefetch.queue('%s')", href)))
	}
	return l.Element.Render(ctx, w)
}

func (*Link) NodeType() html.NodeType {
	return html.NodeElement
}

var Eager = html.Attr("eager", true)

func A(props ...html.I) html.Element {
	return &Link{Element: html.A(props...)}
}

func Outlet(ctx *Context) html.I {
	return ctx.outlet
}

func Head(ctx *Context) html.I {
	return ctx.head
}

func Body(ctx *Context) html.I {
	return ctx.body
}

func Redirect(ctx *Context, to string) html.I {
	ctx.w.Header().Set("Content-Type", "text/html")
	http.Redirect(ctx.w, ctx.r, to, http.StatusTemporaryRedirect)
	return html.Frag()
}

func NotFound(ctx *Context) html.I {
	return NotFoundPage.Page(ctx)
}

func Error(ctx *Context, errpage ErrorPage) html.I {
	ctx.w.WriteHeader(errpage.Status())
	return errpage.Page(ctx)
}

func SetCookie(ctx *Context, cookies ...*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(ctx.w, cookie)
	}
}

type Header interface {
	Key() string
	Value() string
}

type HeaderEntry struct {
	key   string
	value string
}

func (he *HeaderEntry) Key() string {
	return he.key
}

func (he *HeaderEntry) Value() string {
	return he.value
}

func NewHeader(key, value string) *HeaderEntry {
	return &HeaderEntry{key: key, value: value}
}

func SetHeader(ctx *Context, headers ...Header) {
	for _, header := range headers {
		ctx.w.Header().Set(header.Key(), header.Value())
	}
}

type MetadataOG struct {
	Type        string
	URL         string
	Title       string
	Description string
	// Should the absolute path to an asset
	Image string
}

type MetadataTwitter struct {
	Card        string
	URL         string
	Title       string
	Description string
	// Should the absolute path to an asset
	Image string
}

type Metadata struct {
	Base        string
	Title       string
	Description string
	AppName     string
	Authors     []string
	Generator   string
	Keywords    []string
	Referrer    string
	Creator     string
	Publisher   string
	Robots      string
	Manifest    string
	Icons       string
	Language    string
	Alternates  struct {
		Canonical string
		Languages []string
		Media     string
		Types     []string
	}
	OpenGraph *MetadataOG
	Twitter   *MetadataTwitter
	Assets    []string
}

func (m *Metadata) Render(ctx context.Context, w io.Writer) error {
	var title string = "Pacis App"
	if len(m.Title) > 0 {
		title = m.Title
	}

	el := html.Frag(
		html.Title(html.Text(title)),
		html.Meta(html.Charset("UTF-8")),
		html.Meta(html.HttpEquiv("Content-Type"), html.Content("text/html; charset=utf-8")),
		html.Meta(html.Name("title"), html.Content(title)),
		html.Meta(html.Name("viewport"), html.Content("width=device-width, initial-scale=1.0")),

		html.If(
			len(m.Description) > 0,
			html.Meta(html.Name("description"), html.Content(m.Description)),
		),
		html.If(
			len(m.Keywords) > 0,
			html.Meta(html.Name("keywords"), html.Content(strings.Join(m.Keywords, ","))),
		),
		html.If(
			len(m.Robots) > 0,
			html.Meta(html.Name("robots"), html.Content(m.Robots)),
		),
		html.Map(m.Authors, func(author string, i int) html.I {
			return html.Meta(html.Name("author"), html.Content(author))
		}),
		html.IfFn(m.Twitter != nil, func() html.Renderer {
			return html.Frag(
				html.If(
					len(m.Twitter.Card) > 0,
					html.Meta(html.Property("twitter:card"), html.Content(m.Twitter.Card)),
				),
				html.If(
					len(m.Twitter.Title) > 0,
					html.Meta(html.Property("twitter:title"), html.Content(m.Twitter.Title)),
				),
				html.If(
					len(m.Twitter.Description) > 0,
					html.Meta(html.Property("twitter:description"), html.Content(m.Twitter.Description)),
				),
				html.If(
					len(m.Twitter.URL) > 0,
					html.Meta(html.Property("twitter:url"), html.Content(m.Twitter.URL)),
				),
				html.If(
					len(m.Twitter.Image) > 0,
					html.Meta(html.Property("twitter:image"), html.Content(m.Twitter.Image)),
				),
			)
		}),
		html.IfFn(m.OpenGraph != nil, func() html.Renderer {
			return html.Frag(
				html.If(
					len(m.OpenGraph.Type) > 0,
					html.Meta(html.Property("og:type"), html.Content(m.OpenGraph.Type)),
				),
				html.If(
					len(m.OpenGraph.URL) > 0,
					html.Meta(html.Property("og:url"), html.Content(m.OpenGraph.URL)),
				),
				html.If(
					len(m.OpenGraph.Title) > 0,
					html.Meta(html.Property("og:title"), html.Content(m.OpenGraph.Title)),
				),
				html.If(
					len(m.OpenGraph.Description) > 0,
					html.Meta(html.Property("og:description"), html.Content(m.OpenGraph.Description)),
				),
				html.If(
					len(m.OpenGraph.Image) > 0,
					html.Meta(html.Property("og:image"), html.Content(m.OpenGraph.Image)),
				),
			)
		}),
	)

	return el.Render(ctx, w)
}

func (m *Metadata) NodeType() html.NodeType {
	return html.NodeFragment
}
