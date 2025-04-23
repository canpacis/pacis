package components

import (
	"context"
	"io"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	h "github.com/canpacis/pacis/ui/html"
	"github.com/canpacis/pacis/ui/icons"
)

type CodeHighlighter struct {
	Language string
	Source   string
	Props    []h.I
}

func (c *CodeHighlighter) Render(ctx context.Context, w io.Writer) error {
	lexer := lexers.Get(c.Language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	formatter := html.New(
		html.WithClasses(true),
	)
	style := styles.Fallback

	iterator, err := lexer.Tokenise(nil, c.Source)
	if err != nil {
		return err
	}

	return formatter.Format(w, style, iterator)
}

func (c *CodeHighlighter) NodeType() h.NodeType {
	return h.NodeElement
}

// TODO: Add a prop to include the copy button
func Code(src, lang string, props ...h.I) h.Node {
	return h.Div(
		Join(
			props,
			h.Class("relative overflow-x-auto"),

			Button(
				ButtonSizeIcon,
				ButtonVariantGhost,
				h.Class("!size-7 rounded-sm absolute top-2 right-2 md:top-3 md:right-3"),
				On("click", fn("$clipboard", src)),

				icons.Clipboard(h.Class("size-3")),
				h.Span(h.Class("sr-only"), h.Text("Copy to Clipboard")),
			),
			h.Div(
				h.Class("p-3 pr-10 md:p-6 md:pr-12 overflow-auto bg-accent/50 dark:bg-accent/20 rounded-lg"),
				&CodeHighlighter{Language: lang, Source: src, Props: props},
			),
		)...,
	)
}
