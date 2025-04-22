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

	l.Element.AddAttribute(c.On("mouseenter", fmt.Sprintf("$prefetch.queue('%s')", href)))
	l.Element.AddAttribute(c.On("click", fmt.Sprintf("$prefetch.load('%s', $event)", href)))

	_, eager := l.Element.GetAttribute("eager")
	if eager {
		l.Element.RemoveAttribute("eager")
		l.Element.AddAttribute(c.X("intersect", fmt.Sprintf("$prefetch.queue('%s')", href)))
	}
	return l.Element.Render(ctx, w)
}

func (*Link) NodeType() h.NodeType {
	return h.NodeElement
}

var Eager = h.Attr("eager", true)

func A(props ...h.I) h.Element {
	return &Link{Element: h.A(props...)}
}

type Hook struct{}

func (*Hook) Render(ctx context.Context, w io.Writer) error {
	fmt.Println("Hook render")
	pctx, ok := ctx.(*PageContext)
	if ok {
		fmt.Println(pctx)
		// pctx.hookch <- "test hook"
		// pctx.hookch <- "test hook 2"
	}
	return nil
}

func (*Hook) NodeType() h.NodeType {
	return h.NodeFragment
}

// type QueryPaginator struct {
// 	fn      func(int, int) h.I
// 	Page    int `query:"page"`
// 	PerPage int `query:"per_page"`
// }

// func (p *QueryPaginator) Paginate(ctx *PaginatorContext) error {
// 	if err := ctx.Scan(p); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (p *QueryPaginator) Render(ctx context.Context, w io.Writer) error {
// 	return p.fn(p.Page, p.PerPage).Render(ctx, w)
// }

// func PaginateFn(ctx *PageContext, fn func(int, int) h.I) *QueryPaginator {
// 	p := &QueryPaginator{fn: fn}
// 	// ctx.RegisterPaginator(p)
// 	return p
// }

// func Paginate(ctx *PageContext, p Paginator) Paginator {
// 	// ctx.RegisterPaginator(p)
// 	return p
// }
