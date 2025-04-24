package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/canpacis/pacis/pages/internal"
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
	http.Redirect(ctx.w, ctx.r, to, http.StatusFound)
	return html.Frag()
}

func NotFound(ctx *Context) html.I {
	return NotFoundPage.Page(ctx)
}

func Error(ctx *Context, code int, err error) html.I {
	internal.Set(ctx, "error", err)
	internal.Set(ctx, "status", code)

	return ErrorPage.Page(ctx)
}

func Cookie(ctx *Context, cookies ...*http.Cookie) html.I {
	for _, cookie := range cookies {
		http.SetCookie(ctx.w, cookie)
	}
	return html.Frag()
}

type HeaderEntry struct {
	Key   string
	Value string
}

func NewHeader(key, value string) *HeaderEntry {
	return &HeaderEntry{Key: key, Value: value}
}

func Header(ctx *Context, headers ...*HeaderEntry) html.I {
	for _, header := range headers {
		ctx.r.Header.Set(header.Key, header.Value)
	}
	return html.Frag()
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
	Alternates  struct {
		Canonical string
		Languages []string
		Media     string
		Types     []string
	}
	OpenGraph MetadataOG
	Twitter   MetadataTwitter
	Assets    []string
}
