package pages

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"

	c "github.com/canpacis/pacis/ui/components"
	h "github.com/canpacis/pacis/ui/html"
)

func readattr(attr h.Attribute) string {
	var buf bytes.Buffer
	attr.Render(context.Background(), &buf)
	return buf.String()
}

type Link struct {
	h.Element
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
		l.Element.AddAttribute(h.Target("blank"))
		l.Element.AddAttribute(h.Rel("noreferer"))
		return l.Element.Render(ctx, w)
	}

	l.Element.AddAttribute(c.On("mouseenter", fmt.Sprintf("$prefetch.get('%s')", href)))
	l.Element.AddAttribute(c.On("click", fmt.Sprintf("$prefetch.set('%s', $event)", href)))
	return l.Element.Render(ctx, w)
}

func (*Link) NodeType() h.NodeType {
	return h.NodeElement
}

func A(props ...h.I) h.Element {
	return &Link{Element: h.A(props...)}
}
